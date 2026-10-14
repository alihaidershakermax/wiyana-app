package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/auth"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/db"
)

func (s *server) authCompatRoute(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.URL.Path, "/api/v1/auth/") {
		return false
	}
	switch r.URL.Path {
	case "/api/v1/auth/dev-login":
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "طريقة الطلب غير مدعومة")
			return true
		}
		if !s.cfg.AllowDevAuth || strings.EqualFold(s.cfg.Environment, "production") {
			writeError(w, http.StatusForbidden, "تسجيل التطوير غير متاح")
			return true
		}
		var body struct {
			Role string `json:"role"`
		}
		if !readJSON(w, r, &body) {
			return true
		}
		role := strings.ToLower(strings.TrimSpace(body.Role))
		if role == "" {
			role = "owner"
		}
		identity := auth.Identity{ID: "dev-user", TelegramID: "dev-user", Name: "مستخدم الاختبار", Role: role}
		access, err := auth.SignJWT(identity, s.cfg.JWTSecret, 8*time.Hour)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "تعذر إصدار التوكن")
			return true
		}
		refresh, err := auth.SignJWT(identity, s.cfg.JWTSecret, 30*24*time.Hour)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "تعذر إصدار توكن التحديث")
			return true
		}
		writeJSON(w, http.StatusOK, map[string]any{"access_token": access, "refresh_token": refresh, "token_type": "bearer", "user": adminUser(identity)})
		return true
	case "/api/v1/auth/admin-login":
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "طريقة الطلب غير مدعومة")
			return true
		}
		if !s.cfg.AllowDevAuth || strings.EqualFold(s.cfg.Environment, "production") {
			writeError(w, http.StatusForbidden, "تسجيل الإدارة التجريبي غير متاح")
			return true
		}
		var body struct {
			Identifier string `json:"identifier"`
		}
		if !readJSON(w, r, &body) {
			return true
		}
		identity := auth.Identity{ID: strings.TrimSpace(body.Identifier), TelegramID: strings.TrimSpace(body.Identifier), Name: "مسؤول النظام", Role: "owner"}
		if identity.ID == "" {
			identity.ID = "dev-user"
		}
		access, _ := auth.SignJWT(identity, s.cfg.JWTSecret, 8*time.Hour)
		refresh, _ := auth.SignJWT(identity, s.cfg.JWTSecret, 30*24*time.Hour)
		writeJSON(w, http.StatusOK, map[string]any{"access_token": access, "refresh_token": refresh, "token_type": "bearer", "user": adminUser(identity)})
		return true
	case "/api/v1/auth/telegram":
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "طريقة الطلب غير مدعومة")
			return true
		}
		if s.cfg.AllowDevAuth && !strings.EqualFold(s.cfg.Environment, "production") {
			identity := auth.Identity{ID: "dev-user", TelegramID: "dev-user", Name: "مستخدم الاختبار", Role: "owner"}
			access, _ := auth.SignJWT(identity, s.cfg.JWTSecret, 8*time.Hour)
			refresh, _ := auth.SignJWT(identity, s.cfg.JWTSecret, 30*24*time.Hour)
			writeJSON(w, http.StatusOK, map[string]any{"access_token": access, "refresh_token": refresh, "token_type": "bearer", "user": adminUser(identity)})
			return true
		}
		s.loginTelegram(w, r)
		return true
	case "/api/v1/auth/refresh":
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "طريقة الطلب غير مدعومة")
			return true
		}
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if !readJSON(w, r, &body) {
			return true
		}
		request := r.Clone(r.Context())
		request.Header.Set("Authorization", "Bearer "+body.RefreshToken)
		identity, err := auth.Resolve(request, auth.Config{JWTSecret: s.cfg.JWTSecret, BotToken: s.cfg.BotToken, Environment: s.cfg.Environment})
		if err != nil {
			writeError(w, http.StatusUnauthorized, "توكن التحديث غير صالح")
			return true
		}
		access, _ := auth.SignJWT(identity, s.cfg.JWTSecret, 8*time.Hour)
		refresh, _ := auth.SignJWT(identity, s.cfg.JWTSecret, 30*24*time.Hour)
		writeJSON(w, http.StatusOK, map[string]any{"access_token": access, "refresh_token": refresh, "token_type": "bearer"})
		return true
	case "/api/v1/auth/logout":
		if r.Method == http.MethodPost {
			writeJSON(w, http.StatusOK, map[string]any{"success": true})
			return true
		}
		writeError(w, http.StatusMethodNotAllowed, "طريقة الطلب غير مدعومة")
		return true
	case "/api/v1/auth/me":
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "طريقة الطلب غير مدعومة")
			return true
		}
		identity, status := s.identity(w, r)
		if status == 0 {
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "user": adminUser(identity)})
		}
		return true
	default:
		return false
	}
}
func (s *server) compatRoute(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.URL.Path, "/api/v1/") {
		return false
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/")

	// Public endpoints (no auth required — subscriptions removed, all content free)
	if r.Method == http.MethodGet {
		switch {
		case path == "home/dashboard":
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
			return true
		case path == "subjects/" || path == "subjects":
			subjects, _, _ := s.adb.ListSubjects(db.SubjectFilter{Limit: 50})
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": subjects})
			return true
		}
	}

	identity, status := s.identity(w, r)
	if status != 0 {
		return true
	}
	if r.Method == http.MethodGet {
		switch {
		case path == "students/" || path == "students":
			writeJSON(w, http.StatusOK, studentListResponse(r))
			return true
		case path == "teachers/" || path == "teachers":
			writeJSON(w, http.StatusOK, []any{})
			return true
		case path == "audit-logs/" || path == "audit-logs", path == "notifications/" || path == "notifications":
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": []any{}})
			return true
		case path == "settings/" || path == "settings":
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": map[string]any{}})
			return true
		case path == "subscription-codes/" || path == "subscription-codes":
			writeJSON(w, http.StatusOK, paginatedEmptyResponse(r))
			return true
		case path == "panel/subscriptions":
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": []any{}})
			return true
		case path == "analytics/owner" || path == "analytics/teacher":
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": analyticsEmptyResponse(r)})
			return true
		case path == "teacher/dashboard":
			writeJSON(w, http.StatusOK, map[string]any{"my_courses": 0, "my_students": 0, "my_lessons": 0, "total_views": 0, "completion_rate": 0})
			return true
		case path == "teacher/analytics":
			writeJSON(w, http.StatusOK, map[string]any{"period": r.URL.Query().Get("period"), "data": []any{}})
			return true
		case path == "teacher/courses":
			writeJSON(w, http.StatusOK, map[string]any{"courses": []any{}})
			return true
		case path == "users/" || path == "users":
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": []any{}})
			return true
		case path == "courses/" || path == "courses":
			writeJSON(w, http.StatusOK, []any{})
			return true
		case path == "departments/" || path == "departments", path == "favorites/" || path == "favorites", path == "search/" || path == "search", path == "notifications/my-notifications", path == "notifications/unread", path == "subscriptions/my-subscriptions":
			writeJSON(w, http.StatusOK, []any{})
			return true
		case path == "notifications/latest-release":
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": nil})
			return true
		case path == "youtube/status":
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": map[string]any{"configured": false, "is_authenticated": false, "token_file_exists": false, "token_in_db": false, "simulation_mode": false}})
			return true
		case path == "youtube/oauth/authorize":
			writeJSON(w, http.StatusNotImplemented, map[string]any{"success": false, "detail": "ربط YouTube غير مهيأ في البيئة الحالية"})
			return true
		case path == "me" || path == "users/me":
			writeJSON(w, http.StatusOK, homeUser(identity))
			return true
		case path == "dashboard/overview":
			writeJSON(w, http.StatusOK, dashboardOverview())
			return true
		case path == "dashboard/analytics":
			writeJSON(w, http.StatusOK, map[string]any{"period": r.URL.Query().Get("period"), "series": []any{}, "data": []any{}})
			return true
		case path == "dashboard/activity":
			writeJSON(w, http.StatusOK, map[string]any{"activities": []any{}, "data": []any{}})
			return true
		case path == "dashboard/recent-courses":
			writeJSON(w, http.StatusOK, map[string]any{"courses": []any{}})
			return true
		case path == "dashboard/recent-students":
			writeJSON(w, http.StatusOK, map[string]any{"students": []any{}})
			return true
		case strings.HasPrefix(path, "courses/"):
			writeJSON(w, http.StatusNotFound, map[string]any{"detail": "المساق غير موجود"})
			return true
		}
	}
	if r.Method == http.MethodPost || r.Method == http.MethodPatch || r.Method == http.MethodDelete {
		switch {
		case path == "notifications/broadcast", path == "notifications/app-release", path == "settings/" || path == "settings", path == "subscription-codes/batch-generate":
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": []any{}, "message": "تم حفظ العملية محليًا"})
			return true
		case strings.HasPrefix(path, "notifications/") || strings.HasPrefix(path, "subscription-codes/") || strings.HasPrefix(path, "panel/subscriptions/"):
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": map[string]any{}})
			return true
		case path == "students/" || path == "students":
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": []any{}, "meta": map[string]any{"page": 1, "limit": 20, "total": 0, "total_pages": 0}})
			return true
		case strings.HasPrefix(path, "students/") && (strings.HasSuffix(path, "/suspend") || strings.HasSuffix(path, "/activate")):
			writeJSON(w, http.StatusNotFound, map[string]any{"success": false, "detail": "الطالب غير موجود في مصدر البيانات الحالي"})
			return true
		case path == "notifications/read-all", path == "subscriptions/redeem", path == "youtube/oauth/exchange-code", path == "youtube/disconnect":
			writeJSON(w, http.StatusOK, map[string]any{"success": true})
			return true
		case path == "users/me":
			writeJSON(w, http.StatusOK, homeUser(identity))
			return true
		}
	}
	return false
}
func adminUser(identity auth.Identity) map[string]any {
	return map[string]any{"id": identity.ID, "user_id": identity.ID, "telegram_id": identity.TelegramID, "name": identity.Name, "full_name": identity.Name, "role": strings.ToUpper(identity.Role), "is_active": true, "permissions": []string{"courses.read", "students.read", "analytics.read"}}
}
func homeUser(identity auth.Identity) map[string]any {
	return map[string]any{"id": identity.ID, "user_id": identity.ID, "telegram_id": identity.TelegramID, "name": identity.Name, "full_name": identity.Name, "role": identity.Role, "avatar_url": nil}
}
func dashboardOverview() map[string]any {
	return map[string]any{"total_students": 0, "active_students": 0, "total_courses": 0, "published_courses": 0, "total_lessons": 0, "total_views": 0, "total_subscriptions": 0, "revenue": 0, "growth": map[string]any{}, "recent_activity": []any{}}
}
func paginatedEmptyResponse(r *http.Request) map[string]any {
	page := 1
	limit := 20
	if raw := r.URL.Query().Get("page"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	return map[string]any{
		"success": true,
		"data":    []any{},
		"meta": map[string]any{
			"page":        page,
			"limit":       limit,
			"total":       0,
			"total_pages": 0,
		},
	}
}
func analyticsEmptyResponse(r *http.Request) map[string]any {
	return map[string]any{
		"time_range":           r.URL.Query().Get("range"),
		"total_students":       0,
		"active_students":      0,
		"new_students":         0,
		"total_courses":        0,
		"total_lessons":        0,
		"active_subscriptions": 0,
		"views":                0,
		"watch_time_hours":     0,
		"completion_rate":      0,
		"top_courses":          []any{},
	}
}
func studentListResponse(r *http.Request) map[string]any {
	page := 1
	limit := 20
	if raw := r.URL.Query().Get("page"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	return map[string]any{
		"success": true,
		"data":    []any{},
		"meta": map[string]any{
			"page":        page,
			"limit":       limit,
			"total":       0,
			"total_pages": 0,
		},
	}
}
