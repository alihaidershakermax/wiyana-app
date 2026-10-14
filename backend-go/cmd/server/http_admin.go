package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/auth"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/db"
)

// adminPanelRoute serves the admin-panel CRUD API (backed by SQLite).
// It runs before the legacy compat stubs so real handlers win.
// Returns true when the request was handled.
func (s *server) adminPanelRoute(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.URL.Path, "/api/v1/") {
		return false
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1")

	// --- public auth endpoints -------------------------------------------
	switch {
	case path == "/auth/logout" && r.Method == http.MethodPost:
		writeJSON(w, http.StatusOK, map[string]any{"success": true})
		return true
	case path == "/auth/refresh" && r.Method == http.MethodPost:
		s.refreshToken(w, r)
		return true
	case path == "/auth/admin-login" && r.Method == http.MethodPost:
		s.adminLogin(w, r)
		return true
	}

	// --- public content endpoints (no auth — subscriptions removed) --------
	switch {
	case (path == "/subjects" || path == "/subjects/") && r.Method == http.MethodGet:
		if s.adb != nil {
			subjects, _, _ := s.adb.ListSubjects(db.SubjectFilter{Limit: 50})
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": subjects})
		} else {
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": []any{}})
		}
		return true
	case path == "/home/dashboard" && r.Method == http.MethodGet:
		if s.adb != nil {
			courses, _, _ := s.adb.ListCourses(db.CourseFilter{Status: "published", Limit: 10})
			latest := make([]any, 0)
			for _, c := range courses {
				latest = append(latest, map[string]any{
					"id": c.ID, "title": c.Title, "description": c.Description,
					"thumbnail_url": c.ThumbnailURL, "subject_name": c.SubjectName,
					"lessons_count": c.LessonsCount, "created_at": c.CreatedAt,
				})
			}
			writeJSON(w, http.StatusOK, map[string]any{"user": nil, "announcement": nil, "continueLearningCourses": []any{}, "recommendedCourses": latest, "latestCourses": latest, "popularCourses": latest})
		} else {
			writeJSON(w, http.StatusOK, map[string]any{"latestCourses": []any{}})
		}
		return true
	}

	// --- everything below requires an admin/owner/teacher identity ------
	identity, isAdmin := s.admin(w, r)
	if !isAdmin {
		return true // s.admin already wrote the error
	}
	if s.adb == nil {
		writeError(w, http.StatusServiceUnavailable, "قاعدة البيانات غير متاحة")
		return true
	}
	d := s.adb
	actor := identity.ID
	q := r.URL.Query()

	// --- dashboard ---------------------------------------------------------
	switch {
	case path == "/dashboard/overview" && r.Method == http.MethodGet:
		ov, _ := d.Overview()
		writeJSON(w, http.StatusOK, ov)
		return true
	case path == "/dashboard/analytics" && r.Method == http.MethodGet:
		writeJSON(w, http.StatusOK, s.dashboardAnalytics(d, q.Get("period")))
		return true
	case path == "/dashboard/activity" && r.Method == http.MethodGet:
		limit, _ := strconv.Atoi(q.Get("limit"))
		items, _ := d.RecentActivity(limit)
		writeJSON(w, http.StatusOK, map[string]any{"activities": items})
		return true
	case path == "/dashboard/recent-courses" && r.Method == http.MethodGet:
		limit, _ := strconv.Atoi(q.Get("limit"))
		courses, _ := d.RecentCourses(limit)
		writeJSON(w, http.StatusOK, map[string]any{"courses": courses})
		return true
	case path == "/dashboard/recent-students" && r.Method == http.MethodGet:
		limit, _ := strconv.Atoi(q.Get("limit"))
		students, _ := d.RecentStudents(limit)
		writeJSON(w, http.StatusOK, map[string]any{"students": students})
		return true
	}

	// --- teacher self-service ----------------------------------------------
	if strings.HasPrefix(path, "/teacher/") {
		s.teacherRoute(w, r, d, identity, path)
		return true
	}

	// --- youtube (not configured: honest stub) ------------------------------
	switch {
	case path == "/youtube/status" && r.Method == http.MethodGet:
		ok(w, map[string]any{"configured": false, "is_authenticated": false,
			"token_file_exists": false, "token_in_db": false, "simulation_mode": true})
		return true
	case path == "/youtube/oauth/exchange-code" && r.Method == http.MethodPost:
		writeError(w, http.StatusNotImplemented, "ربط يوتيوب غير مفعّل")
		return true
	case path == "/youtube/disconnect" && r.Method == http.MethodDelete:
		ok(w, map[string]any{"success": true})
		return true
	}

	// --- subjects ------------------------------------------------------------
	switch {
	case (path == "/subjects" || path == "/subjects/") && r.Method == http.MethodGet:
		items, total := listSubjects(d, q)
		writeJSON(w, http.StatusOK, paginated(items, total, q))
		return true
	case (path == "/subjects" || path == "/subjects/") && r.Method == http.MethodPost:
		var body struct {
			Name string `json:"name"`
			Code string `json:"code"`
			Description string `json:"description"`
			Department string `json:"department"`
			DepartmentID string `json:"department_id"`
			Stage string `json:"stage"`
			Status string `json:"status"`
		}
		if !readJSON(w, r, &body) {
			return true
		}
		if strings.TrimSpace(body.Name) == "" {
			writeError(w, http.StatusBadRequest, "الاسم مطلوب")
			return true
		}
		sub, err := d.CreateSubject(body.Name, body.Code, body.Description, body.Department, body.DepartmentID, body.Stage, body.Status, actor)
		if err != nil {
			writeError(w, http.StatusBadRequest, "تعذر إنشاء المادة")
			return true
		}
		d.Audit(actor, "create", "subject", sub.ID, body.Name)
		writeJSON(w, http.StatusCreated, sub)
		return true
	}
	if id, rest, ok2 := split2(path, "/subjects/"); ok2 && rest == "" {
		switch r.Method {
		case http.MethodGet:
			sub, err := d.GetSubject(id)
			if err != nil {
				writeError(w, http.StatusNotFound, "المادة غير موجودة")
				return true
			}
			writeJSON(w, http.StatusOK, sub)
			return true
		case http.MethodPatch, http.MethodPut:
			fields := patchFields(w, r)
			if fields == nil {
				return true
			}
			sub, err := d.UpdateSubject(id, fields)
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر تحديث المادة")
				return true
			}
			d.Audit(actor, "update", "subject", id, "")
			writeJSON(w, http.StatusOK, sub)
			return true
		case http.MethodDelete:
			deleted, coursesCount, err := d.DeleteSubject(id)
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر حذف المادة")
				return true
			}
			d.Audit(actor, "delete", "subject", id, "")
			msg := "تم حذف المادة"
			if !deleted {
				msg = "تمت أرشفة المادة لوجود كورسات مرتبطة بها"
			}
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "deleted": deleted,
				"archived": !deleted, "message": msg, "courses_count": coursesCount})
			return true
		}
	}

	// --- courses --------------------------------------------------------------
	switch {
	case (path == "/courses" || path == "/courses/") && r.Method == http.MethodGet:
		items, total := listCourses(d, q)
		writeJSON(w, http.StatusOK, paginated(items, total, q))
		return true
	case (path == "/courses" || path == "/courses/") && r.Method == http.MethodPost:
		var body struct {
			Title string `json:"title"`
			Description string `json:"description"`
			ThumbnailURL string `json:"thumbnail_url"`
			SubjectID string `json:"subject_id"`
			Department string `json:"department"`
			DepartmentID string `json:"department_id"`
			Stage string `json:"stage"`
			TeacherID string `json:"teacher_id"`
			Status string `json:"status"`
			Price float64 `json:"price"`
		}
		if !readJSON(w, r, &body) {
			return true
		}
		if strings.TrimSpace(body.Title) == "" {
			writeError(w, http.StatusBadRequest, "العنوان مطلوب")
			return true
		}
		c, err := d.CreateCourse(body.Title, body.Description, body.ThumbnailURL, body.SubjectID,
			body.Department, body.DepartmentID, body.Stage, body.TeacherID, body.Price, body.Status, actor)
		if err != nil {
			writeError(w, http.StatusBadRequest, "تعذر إنشاء الكورس")
			return true
		}
		d.Audit(actor, "create", "course", c.ID, body.Title)
		writeJSON(w, http.StatusCreated, c)
		return true
	}
	if id, rest, ok2 := split2(path, "/courses/"); ok2 {
		switch {
		case rest == "" && r.Method == http.MethodGet:
			c, err := d.GetCourse(id)
			if err != nil {
				writeError(w, http.StatusNotFound, "الكورس غير موجود")
				return true
			}
			writeJSON(w, http.StatusOK, c)
			return true
		case rest == "" && (r.Method == http.MethodPatch || r.Method == http.MethodPut):
			fields := patchFields(w, r)
			if fields == nil {
				return true
			}
			c, err := d.UpdateCourse(id, fields)
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر تحديث الكورس")
				return true
			}
			d.Audit(actor, "update", "course", id, "")
			writeJSON(w, http.StatusOK, c)
			return true
		case rest == "" && r.Method == http.MethodDelete:
			if err := d.DeleteCourse(id); err != nil {
				writeError(w, http.StatusBadRequest, "تعذر حذف الكورس")
				return true
			}
			d.Audit(actor, "delete", "course", id, "")
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "تم حذف الكورس"})
			return true
		case rest == "/publish" && r.Method == http.MethodPost:
			c, err := d.SetCourseStatus(id, "published")
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر النشر")
				return true
			}
			d.Audit(actor, "publish", "course", id, "")
			writeJSON(w, http.StatusOK, c)
			return true
		case rest == "/archive" && r.Method == http.MethodPost:
			c, err := d.SetCourseStatus(id, "archived")
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر الأرشفة")
				return true
			}
			d.Audit(actor, "archive", "course", id, "")
			writeJSON(w, http.StatusOK, c)
			return true
		case rest == "/teacher" && r.Method == http.MethodPost:
			var body struct {
				TeacherID string `json:"teacher_id"`
			}
			if !readJSON(w, r, &body) {
				return true
			}
			c, err := d.AssignCourseTeacher(id, body.TeacherID)
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر تعيين المدرس")
				return true
			}
			d.Audit(actor, "assign-teacher", "course", id, body.TeacherID)
			writeJSON(w, http.StatusOK, c)
			return true
		case rest == "/teacher" && r.Method == http.MethodDelete:
			c, err := d.AssignCourseTeacher(id, "")
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر إزالة المدرس")
				return true
			}
			writeJSON(w, http.StatusOK, c)
			return true
		case rest == "/teacher" && r.Method == http.MethodGet:
			c, err := d.GetCourse(id)
			if err != nil {
				writeError(w, http.StatusNotFound, "الكورس غير موجود")
				return true
			}
			var t *db.User
			if c.TeacherID != "" {
				t, _ = d.GetUser(c.TeacherID)
			}
			ok(w, t)
			return true
		case rest == "/sections" && r.Method == http.MethodGet:
			sections, err := d.ListSections(id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "تعذر قراءة الأقسام")
				return true
			}
			if sections == nil {
				sections = []db.Section{}
			}
			writeJSON(w, http.StatusOK, sections) // raw array (sectionApi expects this)
			return true
		case rest == "/sections" && r.Method == http.MethodPost:
			var body struct {
				Title string `json:"title"`
				Description string `json:"description"`
				Order int `json:"order"`
			}
			if !readJSON(w, r, &body) {
				return true
			}
			if strings.TrimSpace(body.Title) == "" {
				writeError(w, http.StatusBadRequest, "العنوان مطلوب")
				return true
			}
			sec, err := d.CreateSection(id, body.Title, body.Description, body.Order)
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر إنشاء القسم")
				return true
			}
			d.Audit(actor, "create", "section", sec.ID, body.Title)
			writeJSON(w, http.StatusCreated, sec)
			return true
		case rest == "/sections/reorder" && r.Method == http.MethodPost:
			var body struct {
				Order []string `json:"order"`
			}
			if !readJSON(w, r, &body) {
				return true
			}
			sections, err := d.ReorderSections(id, body.Order)
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر الترتيب")
				return true
			}
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "sections": sections})
			return true
		case rest == "/lessons" && r.Method == http.MethodGet:
			lessons, err := d.ListLessonsByCourse(id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "تعذر قراءة الدروس")
				return true
			}
			ok(w, lessons)
			return true
		case rest == "/lessons" && r.Method == http.MethodPost:
			lesson, ok3 := s.createLessonFromBody(w, r, d, actor, id, "")
			if ok3 {
				writeJSON(w, http.StatusCreated, lesson)
			}
			return true
		case rest == "/students" && r.Method == http.MethodGet:
			students := courseStudents(d, id)
			ok(w, students)
			return true
		}
	}

	// --- sections ---------------------------------------------------------------
	if id, rest, ok2 := split2(path, "/sections/"); ok2 {
		switch {
		case rest == "" && (r.Method == http.MethodPatch || r.Method == http.MethodPut):
			fields := patchFields(w, r)
			if fields == nil {
				return true
			}
			sec, err := d.UpdateSection(id, fields)
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر تحديث القسم")
				return true
			}
			writeJSON(w, http.StatusOK, sec)
			return true
		case rest == "" && r.Method == http.MethodDelete:
			if err := d.DeleteSection(id); err != nil {
				writeError(w, http.StatusBadRequest, "تعذر حذف القسم")
				return true
			}
			d.Audit(actor, "delete", "section", id, "")
			writeJSON(w, http.StatusOK, map[string]any{"success": true})
			return true
		case rest == "/lessons" && r.Method == http.MethodGet:
			lessons, err := d.ListLessonsBySection(id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "تعذر قراءة الدروس")
				return true
			}
			ok(w, lessons)
			return true
		case rest == "/lessons" && r.Method == http.MethodPost:
			sec, err := d.SectionWithCount(id)
			if err != nil {
				writeError(w, http.StatusNotFound, "القسم غير موجود")
				return true
			}
			lesson, ok3 := s.createLessonFromBody(w, r, d, actor, sec.CourseID, id)
			if ok3 {
				writeJSON(w, http.StatusCreated, lesson)
			}
			return true
		case rest == "/lessons/reorder" && r.Method == http.MethodPost:
			var body struct {
				Order []string `json:"order"`
			}
			if !readJSON(w, r, &body) {
				return true
			}
			lessons, err := d.ReorderLessons(id, body.Order)
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر الترتيب")
				return true
			}
			ok(w, lessons)
			return true
		}
	}

	// --- lessons ------------------------------------------------------------------
	if id, rest, ok2 := split2(path, "/lessons/"); ok2 {
		switch {
		case rest == "" && r.Method == http.MethodGet:
			lesson, err := d.GetLesson(id)
			if err != nil {
				writeError(w, http.StatusNotFound, "الدرس غير موجود")
				return true
			}
			ok(w, lesson)
			return true
		case rest == "" && (r.Method == http.MethodPatch || r.Method == http.MethodPut):
			fields := patchFields(w, r)
			if fields == nil {
				return true
			}
			lesson, err := d.UpdateLesson(id, fields)
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر تحديث الدرس")
				return true
			}
			d.Audit(actor, "update", "lesson", id, "")
			ok(w, lesson)
			return true
		case rest == "" && r.Method == http.MethodDelete:
			if err := d.DeleteLesson(id); err != nil {
				writeError(w, http.StatusBadRequest, "تعذر حذف الدرس")
				return true
			}
			d.Audit(actor, "delete", "lesson", id, "")
			writeJSON(w, http.StatusOK, map[string]any{"success": true})
			return true
		case rest == "/publish" && r.Method == http.MethodPost:
			lesson, err := d.SetLessonStatus(id, "published")
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر النشر")
				return true
			}
			d.Audit(actor, "publish", "lesson", id, "")
			ok(w, lesson)
			return true
		case rest == "/archive" && r.Method == http.MethodPost:
			lesson, err := d.SetLessonStatus(id, "archived")
			if err != nil {
				writeError(w, http.StatusBadRequest, "تعذر الأرشفة")
				return true
			}
			d.Audit(actor, "archive", "lesson", id, "")
			ok(w, lesson)
			return true
		case rest == "/player" && r.Method == http.MethodGet:
			player, err := s.lessonPlayer(d, id)
			if err != nil {
				writeError(w, http.StatusNotFound, "الدرس غير موجود")
				return true
			}
			ok(w, player)
			return true
		case rest == "/upload-video" && r.Method == http.MethodPost:
			s.uploadLessonVideo(w, r, d, actor, id)
			return true
		case rest == "/attachments" && r.Method == http.MethodGet:
			items, err := d.ListAttachments(id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "تعذر قراءة المرفقات")
				return true
			}
			ok(w, items)
			return true
		case rest == "/attachments" && r.Method == http.MethodPost:
			s.uploadAttachment(w, r, d, actor, id)
			return true
		case rest == "/progress" && r.Method == http.MethodGet:
			studentID := q.Get("student_id")
			if studentID == "" {
				studentID = identity.ID
			}
			items, _ := d.ListProgressByStudent(studentID)
			ok(w, items)
			return true
		case rest == "/progress" && r.Method == http.MethodPost:
			var body struct {
				StudentID      string `json:"student_id"`
				WatchedSeconds int    `json:"watched_seconds"`
				Progress       float64 `json:"progress"`
				LastPosition   int    `json:"last_position"`
			}
			if !readJSON(w, r, &body) {
				return true
			}
			if body.StudentID == "" {
				body.StudentID = identity.ID
			}
			completed := 0
			if body.Progress >= 0.9 {
				completed = 1
			}
			_, _ = d.SQL.Exec(`INSERT INTO progress(id,student_id,lesson_id,watched_seconds,completed,updated_at)
				VALUES(?,?,?,?,?,?) ON CONFLICT(student_id,lesson_id) DO UPDATE SET watched_seconds=excluded.watched_seconds, completed=excluded.completed, updated_at=excluded.updated_at`,
				db.NewID(), body.StudentID, id, body.WatchedSeconds, completed, db.NowUTC())
			items, _ := d.ListProgressByStudent(body.StudentID)
			for _, p := range items {
				if p.LessonID == id {
					ok(w, p)
					return true
				}
			}
			ok(w, map[string]any{})
			return true
		}
	}

	// --- attachments ---------------------------------------------------------------
	if id, rest, ok2 := split2(path, "/attachments/"); ok2 && rest == "" && r.Method == http.MethodDelete {
		if err := d.DeleteAttachment(id); err != nil {
			writeError(w, http.StatusBadRequest, "تعذر حذف المرفق")
			return true
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true})
		return true
	}

	// --- students -------------------------------------------------------------------
	switch {
	case (path == "/students" || path == "/students/") && r.Method == http.MethodGet:
		items, total := listStudents(d, q)
		writeJSON(w, http.StatusOK, paginated(items, total, q))
		return true
	}
	if id, rest, ok2 := split2(path, "/students/"); ok2 {
		switch {
		case rest == "" && r.Method == http.MethodGet:
			u, err := d.GetUser(id)
			if err != nil {
				writeError(w, http.StatusNotFound, "الطالب غير موجود")
				return true
			}
			ok(w, u)
			return true
		case rest == "" && (r.Method == http.MethodPatch || r.Method == http.MethodPut):
			fields := patchFields(w, r)
			if fields == nil {
				return true
			}
			allowed := map[string]bool{"first_name": true, "last_name": true, "username": true, "phone": true, "bio": true, "photo_url": true}
			var sets []string
			var args []any
			for k, v := range fields {
				if allowed[k] {
					sets = append(sets, k+" = ?")
					args = append(args, v)
				}
			}
			if len(sets) > 0 {
				args = append(args, db.NowUTC(), id)
				_, _ = d.SQL.Exec("UPDATE users SET "+joinComma(sets)+", updated_at = ? WHERE id = ?", args...)
			}
			u, _ := d.GetUser(id)
			ok(w, u)
			return true
		case rest == "/activate" && r.Method == http.MethodPost:
			_ = d.SetUserActive(id, true)
			d.Audit(actor, "activate", "student", id, "")
			u, _ := d.GetUser(id)
			ok(w, u)
			return true
		case rest == "/suspend" && r.Method == http.MethodPost:
			_ = d.SetUserActive(id, false)
			d.Audit(actor, "suspend", "student", id, "")
			u, _ := d.GetUser(id)
			ok(w, u)
			return true
		case rest == "/activity" && r.Method == http.MethodGet:
			items, _ := d.ListProgressByStudent(id)
			ok(w, items)
			return true
		case rest == "/courses" && r.Method == http.MethodGet:
			items := studentCourses(d, id)
			ok(w, items)
			return true
		case rest == "/progress" && r.Method == http.MethodGet:
			items, _ := d.ListProgressByStudent(id)
			ok(w, items)
			return true
		case rest == "/enrollments" && r.Method == http.MethodGet:
			items, _ := d.ListEnrollmentsByStudent(id)
			ok(w, items)
			return true
		case rest == "/enrollments" && r.Method == http.MethodPost:
			var body struct {
				CourseID string `json:"course_id"`
			}
			if !readJSON(w, r, &body) {
				return true
			}
			e, err := d.Enroll(id, body.CourseID, "active")
			if err != nil || body.CourseID == "" {
				writeError(w, http.StatusBadRequest, "تعذر التسجيل")
				return true
			}
			d.Audit(actor, "enroll", "student", id, body.CourseID)
			ok(w, e)
			return true
		}
		if sid, courseID, ok3 := split3(path, "/students/", "/courses/"); ok3 && strings.HasSuffix(path, "/progress") && r.Method == http.MethodGet {
			cp, _ := d.CourseProgress(sid, courseID)
			completedIDs := completedLessonIDs(d, sid, courseID)
			ok(w, map[string]any{
				"student_id": sid, "course_id": courseID,
				"total_lessons": cp.TotalLessons, "completed_lessons": cp.CompletedLessons,
				"progress_percent": cp.Percent, "completed_lesson_ids": completedIDs,
				"last_watched_lesson_id": lastWatched(d, sid, courseID), "updated_at": db.NowUTC(),
			})
			_ = courseID
			return true
		}
	}

	// --- teachers --------------------------------------------------------------------
	switch {
	case (path == "/teachers" || path == "/teachers/") && r.Method == http.MethodGet:
		users, _, _ := d.ListUsers(db.UserFilter{Role: "teacher", Search: q.Get("search"), Limit: 200})
		if users == nil {
			users = []db.User{}
		}
		writeJSON(w, http.StatusOK, users) // raw array
		return true
	case (path == "/teachers" || path == "/teachers/") && r.Method == http.MethodPost:
		var body struct {
			Name       string `json:"name"`
			Username   string `json:"username"`
			TelegramID int64  `json:"telegram_id"`
			Bio        string `json:"bio"`
		}
		if !readJSON(w, r, &body) {
			return true
		}
		first, last := splitName(body.Name)
		u, err := d.UpsertByTelegram(body.TelegramID, body.Username, first, last, "", "teacher")
		if err != nil {
			writeError(w, http.StatusBadRequest, "تعذر إنشاء المدرس")
			return true
		}
		if body.Bio != "" {
			_, _ = d.SQL.Exec("UPDATE users SET bio=? WHERE id=?", body.Bio, u.ID)
			u.Bio = body.Bio
		}
		// ensure role is teacher even if the telegram id existed before
		_ = d.SetUserRole(u.ID, "teacher")
		u.Role = "teacher"
		d.Audit(actor, "create", "teacher", u.ID, body.Name)
		writeJSON(w, http.StatusCreated, u)
		return true
	}
	if id, rest, ok2 := split2(path, "/teachers/"); ok2 {
		switch {
		case rest == "" && r.Method == http.MethodGet:
			u, err := d.GetUser(id)
			if err != nil {
				writeError(w, http.StatusNotFound, "المدرس غير موجود")
				return true
			}
			writeJSON(w, http.StatusOK, u)
			return true
		case rest == "/status" && (r.Method == http.MethodPatch || r.Method == http.MethodPut):
			var body struct {
				IsActive bool `json:"is_active"`
			}
			if !readJSON(w, r, &body) {
				return true
			}
			_ = d.SetUserActive(id, body.IsActive)
			u, _ := d.GetUser(id)
			writeJSON(w, http.StatusOK, u)
			return true
		}
	}

	return false
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// ok writes the standard {success:true, data:...} envelope.
func ok(w http.ResponseWriter, data any) {
	if data == nil {
		data = map[string]any{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": data})
}

// split2 splits "/base/{id}{rest}" -> (id, rest).
func split2(path, base string) (id, rest string, ok bool) {
	if !strings.HasPrefix(path, base) {
		return "", "", false
	}
	rem := strings.TrimPrefix(path, base)
	if rem == "" || strings.Contains(rem, "?") {
		return "", "", false
	}
	if i := strings.Index(rem, "/"); i >= 0 {
		return rem[:i], rem[i:], true
	}
	return rem, "", true
}

// split3 splits "/students/{sid}/courses/{rest}" -> (sid, rest).
func split3(path, p1, p2 string) (a, b string, ok bool) {
	if !strings.HasPrefix(path, p1) {
		return "", "", false
	}
	rem := strings.TrimPrefix(path, p1)
	i := strings.Index(rem, p2)
	if i < 0 {
		return "", "", false
	}
	return rem[:i], strings.TrimPrefix(rem[i:], p2), true
}

func joinComma(s []string) string { return strings.Join(s, ", ") }

// patchFields decodes a JSON object body into a map (partial updates allowed).
func patchFields(w http.ResponseWriter, r *http.Request) map[string]any {
	var fields map[string]any
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&fields); err != nil {
		writeError(w, http.StatusBadRequest, "صيغة الطلب غير صالحة")
		return nil
	}
	return fields
}

func pageLimit(q url.Values) (page, limit int) {
	page, limit = 1, 20
	if v, err := strconv.Atoi(q.Get("page")); err == nil && v > 0 {
		page = v
	}
	if v, err := strconv.Atoi(q.Get("limit")); err == nil && v > 0 && v <= 200 {
		limit = v
	}
	return page, limit
}

func paginated(items any, total int, q url.Values) map[string]any {
	page, limit := pageLimit(q)
	totalPages := (total + limit - 1) / limit
	return map[string]any{"success": true, "data": items,
		"meta": map[string]any{"total": total, "page": page, "limit": limit, "total_pages": totalPages}}
}

func listSubjects(d *db.DB, q url.Values) ([]db.Subject, int) {
	page, limit := pageLimit(q)
	items, total, _ := d.ListSubjects(db.SubjectFilter{Search: q.Get("search"), Department: q.Get("department"),
		Stage: q.Get("stage"), Status: q.Get("status"), Page: page, Limit: limit})
	if items == nil {
		items = []db.Subject{}
	}
	return items, total
}

func listCourses(d *db.DB, q url.Values) ([]db.Course, int) {
	page, limit := pageLimit(q)
	items, total, _ := d.ListCourses(db.CourseFilter{Search: q.Get("search"), SubjectID: q.Get("subject_id"),
		TeacherID: q.Get("teacher_id"), Department: q.Get("department"), Stage: q.Get("stage"),
		Status: q.Get("status"), Page: page, Limit: limit})
	if items == nil {
		items = []db.Course{}
	}
	return items, total
}

func listStudents(d *db.DB, q url.Values) ([]db.User, int) {
	page, limit := pageLimit(q)
	role := q.Get("role")
	if role == "" {
		role = "student"
	}
	items, total, _ := d.ListUsers(db.UserFilter{Role: role, Search: q.Get("search"), Page: page, Limit: limit})
	if items == nil {
		items = []db.User{}
	}
	return items, total
}

func courseStudents(d *db.DB, courseID string) []any {
	enrolls, _ := d.ListEnrollmentsByCourse(courseID)
	out := []any{}
	for _, e := range enrolls {
		u, err := d.GetUser(e.StudentID)
		if err != nil {
			continue
		}
		out = append(out, u)
	}
	return out
}

func studentCourses(d *db.DB, studentID string) []any {
	enrolls, _ := d.ListEnrollmentsByStudent(studentID)
	out := []any{}
	for _, e := range enrolls {
		c, err := d.GetCourse(e.CourseID)
		if err != nil {
			continue
		}
		progress, _ := d.ListProgressByStudent(studentID)
		out = append(out, map[string]any{
			"course":     map[string]any{"id": c.ID, "title": c.Title, "teacher_id": c.TeacherID, "status": c.Status},
			"enrollment": e,
			"progress":   progress,
		})
	}
	return out
}

func completedLessonIDs(d *db.DB, studentID, courseID string) []string {
	rows, _ := d.SQL.Query(`SELECT p.lesson_id FROM progress p JOIN lessons l ON l.id=p.lesson_id
		WHERE p.student_id=? AND l.course_id=? AND p.completed=1`, studentID, courseID)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	out := []string{}
	if rows == nil {
		return out
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			out = append(out, id)
		}
	}
	return out
}

func lastWatched(d *db.DB, studentID, courseID string) string {
	var id string
	_ = d.SQL.QueryRow(`SELECT p.lesson_id FROM progress p JOIN lessons l ON l.id=p.lesson_id
		WHERE p.student_id=? AND l.course_id=? ORDER BY p.updated_at DESC LIMIT 1`, studentID, courseID).Scan(&id)
	return id
}

func splitName(name string) (first, last string) {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

// createLessonFromBody handles POST /courses/{id}/lessons and /sections/{id}/lessons.
func (s *server) createLessonFromBody(w http.ResponseWriter, r *http.Request, d *db.DB, actor, courseID, sectionID string) (*db.Lesson, bool) {
	var body struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		SectionID      string `json:"section_id"`
		Thumbnail      string `json:"thumbnail"`
		YoutubeVideoID string `json:"youtube_video_id"`
		YoutubeURL     string `json:"youtube_url"`
		Duration       int    `json:"duration"`
		Order          int    `json:"order"`
		Status         string `json:"status"`
		IsFree         bool   `json:"is_free"`
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "صيغة الطلب غير صالحة")
		return nil, false
	}
	if strings.TrimSpace(body.Title) == "" {
		writeError(w, http.StatusBadRequest, "العنوان مطلوب")
		return nil, false
	}
	if body.SectionID != "" {
		sectionID = body.SectionID
	}
	lesson, err := d.CreateLesson(courseID, sectionID, body.Title, body.Description, body.Thumbnail,
		body.YoutubeVideoID, body.YoutubeURL, body.Duration, body.Order, body.Status, body.IsFree, actor)
	if err != nil {
		writeError(w, http.StatusBadRequest, "تعذر إنشاء الدرس")
		return nil, false
	}
	d.Audit(actor, "create", "lesson", lesson.ID, body.Title)
	return lesson, true
}

// lessonPlayer builds the player payload for a lesson.
func (s *server) lessonPlayer(d *db.DB, lessonID string) (map[string]any, error) {
	lesson, err := d.GetLesson(lessonID)
	if err != nil {
		return nil, err
	}
	course, _ := d.GetCourse(lesson.CourseID)
	var section *db.Section
	if lesson.SectionID != "" {
		section, _ = d.SectionWithCount(lesson.SectionID)
	}
	attachments, _ := d.ListAttachments(lessonID)
	// previous / next within the same section (fallback: course)
	var prev, next *db.Lesson
	siblings, _ := d.ListLessonsBySection(lesson.SectionID)
	if len(siblings) == 0 {
		siblings, _ = d.ListLessonsByCourse(lesson.CourseID)
	}
	for i, l := range siblings {
		if l.ID == lessonID {
			if i > 0 {
				c := siblings[i-1]
				prev = &c
			}
			if i+1 < len(siblings) {
				c := siblings[i+1]
				next = &c
			}
			break
		}
	}
	courseInfo := map[string]any{"id": "", "title": ""}
	if course != nil {
		courseInfo = map[string]any{"id": course.ID, "title": course.Title}
	}
	return map[string]any{
		"lesson": lesson, "course": courseInfo, "section": section,
		"attachments": attachments, "previous_lesson": prev, "next_lesson": next,
	}, nil
}

// uploadLessonVideo handles multipart video upload for a lesson.
func (s *server) uploadLessonVideo(w http.ResponseWriter, r *http.Request, d *db.DB, actor, lessonID string) {
	if _, err := d.GetLesson(lessonID); err != nil {
		writeError(w, http.StatusNotFound, "الدرس غير موجود")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 2<<30) // 2GB
	if err := r.ParseMultipartForm(2 << 30); err != nil {
		writeError(w, http.StatusBadRequest, "تعذر قراءة الملف")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "الملف مطلوب")
		return
	}
	defer file.Close()
	name := sanitizeFilename(header.Filename)
	if name == "" {
		name = "video.mp4"
	}
	dir := filepath.Join(s.cfg.UploadDir, "videos")
	_ = os.MkdirAll(dir, 0o755)
	stored := db.NewID() + "_" + name
	dst, err := os.Create(filepath.Join(dir, stored))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "تعذر حفظ الملف")
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		_ = dst.Close()
		writeError(w, http.StatusInternalServerError, "تعذر حفظ الملف")
		return
	}
	_ = dst.Close()
	rel := filepath.Join("videos", stored)
	lesson, err := d.SetLessonVideo(lessonID, rel)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "تعذر ربط الفيديو بالدرس")
		return
	}
	d.Audit(actor, "upload-video", "lesson", lessonID, name)
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": lesson, "youtube": map[string]any{}})
}

// uploadAttachment handles file or link attachments for a lesson.
func (s *server) uploadAttachment(w http.ResponseWriter, r *http.Request, d *db.DB, actor, lessonID string) {
	if _, err := d.GetLesson(lessonID); err != nil {
		writeError(w, http.StatusNotFound, "الدرس غير موجود")
		return
	}
	ct := r.Header.Get("Content-Type")
	var name, urlStr, storageKey, mimeType string
	var size int64
	if strings.HasPrefix(ct, "multipart/form-data") {
		r.Body = http.MaxBytesReader(w, r.Body, 100<<20)
		if err := r.ParseMultipartForm(100 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "تعذر قراءة الطلب")
			return
		}
		name = r.FormValue("name")
		urlStr = r.FormValue("url")
		if fh, header, err := r.FormFile("file"); err == nil {
			defer fh.Close()
			if name == "" {
				name = header.Filename
			}
			mimeType = header.Header.Get("Content-Type")
			size = header.Size
			dir := filepath.Join(s.cfg.UploadDir, "attachments")
			_ = os.MkdirAll(dir, 0o755)
			stored := db.NewID() + "_" + sanitizeFilename(header.Filename)
			dst, err := os.Create(filepath.Join(dir, stored))
			if err != nil {
				writeError(w, http.StatusInternalServerError, "تعذر حفظ الملف")
				return
			}
			if _, err := io.Copy(dst, fh); err != nil {
				_ = dst.Close()
				writeError(w, http.StatusInternalServerError, "تعذر حفظ الملف")
				return
			}
			_ = dst.Close()
			storageKey = filepath.Join("attachments", stored)
		}
	} else {
		var body struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}
		if !readJSON(w, r, &body) {
			return
		}
		name, urlStr = body.Name, body.URL
	}
	if name == "" {
		name = urlStr
	}
	if name == "" {
		writeError(w, http.StatusBadRequest, "الاسم أو الرابط مطلوب")
		return
	}
	typ := "file"
	if urlStr != "" && storageKey == "" {
		typ = "link"
	}
	a, err := d.CreateAttachment(lessonID, name, typ, urlStr, storageKey, mimeType, size, actor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "تعذر إنشاء المرفق")
		return
	}
	ok(w, a)
}

func sanitizeFilename(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, " ", "_")
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" || out == "." {
		return "file"
	}
	return out
}

// dashboardAnalytics builds simple per-bucket counts.
func (s *server) dashboardAnalytics(d *db.DB, period string) map[string]any {
	days := 7
	switch period {
	case "day":
		days = 1
	case "30d":
		days = 30
	case "3m":
		days = 90
	case "1y":
		days = 365
	}
	type point struct {
		Label         string `json:"label"`
		Students      int    `json:"students"`
		Courses       int    `json:"courses"`
		Views         int    `json:"views"`
		Subscriptions int    `json:"subscriptions"`
	}
	points := []point{}
	now := time.Now()
	step := days
	if days > 31 {
		step = days / 12
		if step < 1 {
			step = 1
		}
	}
	for i := step; i >= 1; i-- {
		end := now.AddDate(0, 0, -((i - 1) * step))
		start := end.AddDate(0, 0, -step)
		var students, courses int
		_ = d.SQL.QueryRow("SELECT COUNT(*) FROM users WHERE role='student' AND created_at >= ? AND created_at < ?",
			start.Format(time.RFC3339), end.Format(time.RFC3339)).Scan(&students)
		_ = d.SQL.QueryRow("SELECT COUNT(*) FROM courses WHERE created_at >= ? AND created_at < ?",
			start.Format(time.RFC3339), end.Format(time.RFC3339)).Scan(&courses)
		points = append(points, point{Label: start.Format("2006-01-02"), Students: students, Courses: courses})
	}
	return map[string]any{"period": period, "data": points}
}

// teacherRoute serves /teacher/* for teacher-role users.
func (s *server) teacherRoute(w http.ResponseWriter, r *http.Request, d *db.DB, identity auth.Identity, path string) {
	// map telegram identity -> user row for teacher_id matching
	var teacherID string
	_ = d.SQL.QueryRow("SELECT id FROM users WHERE CAST(telegram_id AS TEXT)=? OR id=? LIMIT 1",
		identity.TelegramID, identity.ID).Scan(&teacherID)
	switch {
	case path == "/teacher/dashboard" && r.Method == http.MethodGet:
		var courses, students, lessons, views int
		_ = d.SQL.QueryRow("SELECT COUNT(*) FROM courses WHERE teacher_id=?", teacherID).Scan(&courses)
		_ = d.SQL.QueryRow("SELECT COUNT(*) FROM enrollments e JOIN courses c ON c.id=e.course_id WHERE c.teacher_id=? AND e.status='active'", teacherID).Scan(&students)
		_ = d.SQL.QueryRow("SELECT COUNT(*) FROM lessons l JOIN courses c ON c.id=l.course_id WHERE c.teacher_id=?", teacherID).Scan(&lessons)
		_ = d.SQL.QueryRow("SELECT COALESCE(SUM(l.views),0) FROM lessons l JOIN courses c ON c.id=l.course_id WHERE c.teacher_id=?", teacherID).Scan(&views)
		writeJSON(w, http.StatusOK, map[string]any{"my_courses": courses, "my_students": students,
			"my_lessons": lessons, "total_views": views, "completion_rate": 0})
	case path == "/teacher/analytics" && r.Method == http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{"period": r.URL.Query().Get("period"), "data": []any{}})
	case path == "/teacher/courses" && r.Method == http.MethodGet:
		items, _, _ := d.ListCourses(db.CourseFilter{TeacherID: teacherID, Limit: 200})
		if items == nil {
			items = []db.Course{}
		}
		writeJSON(w, http.StatusOK, map[string]any{"courses": items})
	default:
		writeError(w, http.StatusNotFound, "المسار غير موجود")
	}
}

// refreshToken issues a new access token from a valid refresh token.
func (s *server) refreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	identity, err := auth.ParseJWT(body.RefreshToken, s.cfg.JWTSecret)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "توكن التحديث غير صالح")
		return
	}
	access, err := auth.SignJWT(identity, s.cfg.JWTSecret, 8*time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "تعذر إصدار التوكن")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"access_token": access, "token_type": "bearer", "user": adminUser(identity)})
}

// adminLogin signs in with a user identifier (telegram id / username / id).
// Only users with an admin-capable role and active accounts may sign in.
func (s *server) adminLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Identifier string `json:"identifier"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	id := strings.TrimSpace(body.Identifier)
	if id == "" || s.adb == nil {
		writeError(w, http.StatusUnauthorized, "بيانات الدخول غير صالحة")
		return
	}
	var u db.User
	err := dbScanUser(s.adb, &u, "SELECT "+db.UserCols+" FROM users WHERE (CAST(telegram_id AS TEXT)=? OR username=? OR id=?) AND is_active=1 LIMIT 1", id, id, id)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "الحساب غير موجود أو غير مفعّل")
		return
	}
	role := strings.ToLower(u.Role)
	if role != "admin" && role != "owner" && role != "teacher" {
		writeError(w, http.StatusForbidden, "لا تملك صلاحية الإدارة")
		return
	}
	identity := auth.Identity{ID: u.ID, TelegramID: strconv.FormatInt(u.TelegramID, 10),
		Name: strings.TrimSpace(u.FirstName + " " + u.LastName), Role: role}
	if identity.Name == "" {
		identity.Name = u.Username
	}
	access, err := auth.SignJWT(identity, s.cfg.JWTSecret, 8*time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "تعذر إصدار التوكن")
		return
	}
	refresh, err := auth.SignJWT(identity, s.cfg.JWTSecret, 30*24*time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "تعذر إصدار توكن التحديث")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"access_token": access, "refresh_token": refresh,
		"token_type": "bearer", "user": adminUser(identity)})
}

// dbScanUser is a small adapter so handlers can scan db.User rows.
func dbScanUser(d *db.DB, u *db.User, query string, args ...any) error {
	row := d.SQL.QueryRow(query, args...)
	var active int
	err := row.Scan(&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.PhotoURL, &u.Role, &active, &u.Phone, &u.Bio, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}
	u.IsActive = active != 0
	return nil
}
