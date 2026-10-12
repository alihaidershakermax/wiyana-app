package db

import (
	"database/sql"
	"strings"
)

// ---------------------------------------------------------------------------
// Model structs (JSON shapes match the admin frontend types)
// ---------------------------------------------------------------------------

type User struct {
	ID         string `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	PhotoURL   string `json:"photo_url"`
	Role       string `json:"role"`
	IsActive   bool   `json:"is_active"`
	Phone      string `json:"phone"`
	Bio        string `json:"bio"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type Subject struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Description  string `json:"description"`
	Department   string `json:"department"`
	DepartmentID string `json:"department_id"`
	Stage        string `json:"stage"`
	Status       string `json:"status"`
	CreatedBy    string `json:"created_by"`
	CoursesCount int    `json:"courses_count"`
	StudentsCount int   `json:"students_count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type Course struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	ThumbnailURL string  `json:"thumbnail_url"`
	SubjectID    string  `json:"subject_id"`
	SubjectName  string  `json:"subject_name"`
	Department   string  `json:"department"`
	DepartmentID string  `json:"department_id"`
	Stage        string  `json:"stage"`
	TeacherID    string  `json:"teacher_id"`
	TeacherName  string  `json:"teacher_name"`
	Price        float64 `json:"price"`
	Status       string  `json:"status"`
	IsPublished  bool    `json:"is_published"`
	StudentsCount int    `json:"students_count"`
	LessonsCount  int    `json:"lessons_count"`
	PublishedAt  string  `json:"published_at"`
	CreatedBy    string  `json:"created_by"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type Section struct {
	ID          string `json:"id"`
	CourseID    string `json:"course_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	LessonsCount int   `json:"lessons_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Lesson struct {
	ID            string `json:"id"`
	CourseID      string `json:"course_id"`
	SectionID     string `json:"section_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Thumbnail     string `json:"thumbnail"`
	FilePath      string `json:"file_path"`
	YoutubeVideoID string `json:"youtube_video_id"`
	YoutubeURL    string `json:"youtube_url"`
	Duration      int    `json:"duration"`
	Order         int    `json:"order"`
	Status        string `json:"status"`
	IsFree        bool   `json:"is_free"`
	Views         int    `json:"views"`
	CreatedBy     string `json:"created_by"`
	PublishedAt   string `json:"published_at"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type Attachment struct {
	ID         string `json:"id"`
	LessonID   string `json:"lesson_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Size       int64  `json:"size"`
	URL        string `json:"url"`
	StorageKey string `json:"storage_key"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
}

type Enrollment struct {
	ID        string `json:"id"`
	StudentID string `json:"student_id"`
	CourseID  string `json:"course_id"`
	CourseTitle string `json:"course_title,omitempty"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

type ActivityItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func scanUser(u *User, row *sql.Row) error {
	var active int
	err := row.Scan(&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.PhotoURL, &u.Role, &active, &u.Phone, &u.Bio, &u.CreatedAt, &u.UpdatedAt)
	u.IsActive = active != 0
	return err
}

func scanUserRows(u *User, rows *sql.Rows) error {
	var active int
	err := rows.Scan(&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.PhotoURL, &u.Role, &active, &u.Phone, &u.Bio, &u.CreatedAt, &u.UpdatedAt)
	u.IsActive = active != 0
	return err
}

const userCols = "id, telegram_id, username, first_name, last_name, photo_url, role, is_active, phone, bio, created_at, updated_at"

// UserCols is the exported column list for handlers outside this package.
const UserCols = userCols

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

// UpsertByTelegram creates or updates a user from Telegram login data,
// preserving an existing role (so a promoted owner is never demoted).
func (d *DB) UpsertByTelegram(telegramID int64, username, firstName, lastName, photoURL, role string) (*User, error) {
	now := NowUTC()
	var existing User
	err := scanUser(&existing, d.SQL.QueryRow("SELECT "+userCols+" FROM users WHERE telegram_id = ?", telegramID))
	if err == nil {
		_, err = d.SQL.Exec(`UPDATE users SET username=?, first_name=?, last_name=?, photo_url=?, updated_at=? WHERE id=?`,
			username, firstName, lastName, photoURL, now, existing.ID)
		if err != nil {
			return nil, err
		}
		existing.Username, existing.FirstName, existing.LastName, existing.PhotoURL, existing.UpdatedAt =
			username, firstName, lastName, photoURL, now
		return &existing, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	u := &User{ID: NewID(), TelegramID: telegramID, Username: username, FirstName: firstName,
		LastName: lastName, PhotoURL: photoURL, Role: role, IsActive: true, CreatedAt: now, UpdatedAt: now}
	_, err = d.SQL.Exec(`INSERT INTO users(id, telegram_id, username, first_name, last_name, photo_url, role, is_active, created_at, updated_at)
		VALUES(?,?,?,?,?,?,?,?,?,?)`, u.ID, u.TelegramID, u.Username, u.FirstName, u.LastName, u.PhotoURL, u.Role, 1, now, now)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (d *DB) GetUser(id string) (*User, error) {
	var u User
	if err := scanUser(&u, d.SQL.QueryRow("SELECT "+userCols+" FROM users WHERE id = ?", id)); err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *DB) SetUserRole(id, role string) error {
	_, err := d.SQL.Exec("UPDATE users SET role=?, updated_at=? WHERE id=?", role, NowUTC(), id)
	return err
}

func (d *DB) SetUserActive(id string, active bool) error {
	v := 0
	if active {
		v = 1
	}
	_, err := d.SQL.Exec("UPDATE users SET is_active=?, updated_at=? WHERE id=?", v, NowUTC(), id)
	return err
}

type UserFilter struct {
	Role   string
	Search string
	Page   int
	Limit  int
}

func (d *DB) ListUsers(f UserFilter) ([]User, int, error) {
	var conds []string
	var args []any
	if f.Role != "" {
		conds = append(conds, "role = ?")
		args = append(args, f.Role)
	}
	if f.Search != "" {
		conds = append(conds, "(username LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR CAST(telegram_id AS TEXT) LIKE ?)")
		q := "%" + f.Search + "%"
		args = append(args, q, q, q, q)
	}
	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}
	var total int
	if err := d.SQL.QueryRow("SELECT COUNT(*) FROM users "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	limit := f.Limit
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	page := f.Page
	if page < 1 {
		page = 1
	}
	rows, err := d.SQL.Query("SELECT "+userCols+" FROM users "+where+" ORDER BY created_at DESC LIMIT ? OFFSET ?",
		append(args, limit, (page-1)*limit)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		var u User
		if err := scanUserRows(&u, rows); err != nil {
			return nil, 0, err
		}
		out = append(out, u)
	}
	return out, total, rows.Err()
}

// ---------------------------------------------------------------------------
// Subjects
// ---------------------------------------------------------------------------

func (d *DB) CreateSubject(name, code, description, department, departmentID, stage, status, createdBy string) (*Subject, error) {
	if status == "" {
		status = "active"
	}
	now := NowUTC()
	s := &Subject{ID: NewID(), Name: name, Code: code, Description: description, Department: department,
		DepartmentID: departmentID, Stage: stage, Status: status, CreatedBy: createdBy, CreatedAt: now, UpdatedAt: now}
	_, err := d.SQL.Exec(`INSERT INTO subjects(id,name,code,description,department,department_id,stage,status,created_by,created_at,updated_at)
		VALUES(?,?,?,?,?,?,?,?,?,?,?)`, s.ID, s.Name, s.Code, s.Description, s.Department, s.DepartmentID, s.Stage, s.Status, s.CreatedBy, now, now)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (d *DB) GetSubject(id string) (*Subject, error) {
	var s Subject
	err := d.SQL.QueryRow(`SELECT id,name,code,description,department,department_id,stage,status,created_by,created_at,updated_at
		FROM subjects WHERE id=?`, id).Scan(&s.ID, &s.Name, &s.Code, &s.Description, &s.Department, &s.DepartmentID, &s.Stage, &s.Status, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM courses WHERE subject_id=?", id).Scan(&s.CoursesCount)
	return &s, nil
}

type SubjectFilter struct {
	Search     string
	Department string
	Stage      string
	Status     string
	Page       int
	Limit      int
}

func (d *DB) ListSubjects(f SubjectFilter) ([]Subject, int, error) {
	var conds []string
	var args []any
	if f.Search != "" {
		conds = append(conds, "(name LIKE ? OR code LIKE ?)")
		q := "%" + f.Search + "%"
		args = append(args, q, q)
	}
	if f.Department != "" {
		conds = append(conds, "department = ?")
		args = append(args, f.Department)
	}
	if f.Stage != "" {
		conds = append(conds, "stage = ?")
		args = append(args, f.Stage)
	}
	if f.Status != "" {
		conds = append(conds, "status = ?")
		args = append(args, f.Status)
	}
	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}
	var total int
	if err := d.SQL.QueryRow("SELECT COUNT(*) FROM subjects "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	limit := f.Limit
	if limit <= 0 || limit > 200 {
		limit = 20
	}
	page := f.Page
	if page < 1 {
		page = 1
	}
	rows, err := d.SQL.Query(`SELECT id,name,code,description,department,department_id,stage,status,created_by,created_at,updated_at
		FROM subjects `+where+` ORDER BY created_at DESC LIMIT ? OFFSET ?`, append(args, limit, (page-1)*limit)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []Subject
	for rows.Next() {
		var s Subject
		if err := rows.Scan(&s.ID, &s.Name, &s.Code, &s.Description, &s.Department, &s.DepartmentID, &s.Stage, &s.Status, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt); err != nil {
			rows.Close()
			return nil, 0, err
		}
		out = append(out, s)
	}
	qerr := rows.Err()
	rows.Close()
	if qerr != nil {
		return nil, 0, qerr
	}
	for i := range out {
		_ = d.SQL.QueryRow("SELECT COUNT(*) FROM courses WHERE subject_id=?", out[i].ID).Scan(&out[i].CoursesCount)
	}
	return out, total, nil
}

func (d *DB) UpdateSubject(id string, fields map[string]any) (*Subject, error) {
	allowed := map[string]bool{"name": true, "code": true, "description": true, "department": true, "department_id": true, "stage": true, "status": true}
	var sets []string
	var args []any
	for k, v := range fields {
		if allowed[k] {
			sets = append(sets, k+" = ?")
			args = append(args, v)
		}
	}
	if len(sets) == 0 {
		return d.GetSubject(id)
	}
	sets = append(sets, "updated_at = ?")
	args = append(args, NowUTC(), id)
	_, err := d.SQL.Exec("UPDATE subjects SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...)
	if err != nil {
		return nil, err
	}
	return d.GetSubject(id)
}

// DeleteSubject archives the subject when courses reference it, deletes otherwise.
func (d *DB) DeleteSubject(id string) (deleted bool, coursesCount int, err error) {
	if err := d.SQL.QueryRow("SELECT COUNT(*) FROM courses WHERE subject_id=?", id).Scan(&coursesCount); err != nil {
		return false, 0, err
	}
	if coursesCount > 0 {
		_, err = d.SQL.Exec("UPDATE subjects SET status='archived', updated_at=? WHERE id=?", NowUTC(), id)
		return false, coursesCount, err
	}
	_, err = d.SQL.Exec("DELETE FROM subjects WHERE id=?", id)
	return err == nil, coursesCount, err
}

// ---------------------------------------------------------------------------
// Courses
// ---------------------------------------------------------------------------

const courseSelect = `SELECT c.id, c.title, c.description, c.thumbnail_url, c.subject_id,
	COALESCE(s.name,''), c.department, c.department_id, c.stage, c.teacher_id,
	COALESCE(u.first_name || ' ' || u.last_name, ''), c.price, c.status, c.is_published,
	c.published_at, c.created_by, c.created_at, c.updated_at
	FROM courses c LEFT JOIN subjects s ON s.id = c.subject_id LEFT JOIN users u ON u.id = c.teacher_id`

func scanCourse(c *Course, rows *sql.Rows) error {
	var published int
	var teacherName string
	err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.ThumbnailURL, &c.SubjectID,
		&c.SubjectName, &c.Department, &c.DepartmentID, &c.Stage, &c.TeacherID,
		&teacherName, &c.Price, &c.Status, &published,
		&c.PublishedAt, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	c.IsPublished = published != 0
	c.TeacherName = strings.TrimSpace(teacherName)
	return err
}

func (d *DB) courseCounts(c *Course) {
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM enrollments WHERE course_id=? AND status='active'", c.ID).Scan(&c.StudentsCount)
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM lessons WHERE course_id=?", c.ID).Scan(&c.LessonsCount)
}

func (d *DB) CreateCourse(title, description, thumbnailURL, subjectID, department, departmentID, stage, teacherID string, price float64, status, createdBy string) (*Course, error) {
	if status == "" {
		status = "draft"
	}
	now := NowUTC()
	published := 0
	var publishedAt string
	if status == "published" {
		published, publishedAt = 1, now
	}
	id := NewID()
	_, err := d.SQL.Exec(`INSERT INTO courses(id,title,description,thumbnail_url,subject_id,department,department_id,stage,teacher_id,price,status,is_published,published_at,created_by,created_at,updated_at)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, id, title, description, thumbnailURL, subjectID, department, departmentID, stage, teacherID, price, status, published, publishedAt, createdBy, now, now)
	if err != nil {
		return nil, err
	}
	return d.GetCourse(id)
}

func (d *DB) GetCourse(id string) (*Course, error) {
	rows, err := d.SQL.Query(courseSelect+" WHERE c.id = ?", id)
	if err != nil {
		return nil, err
	}
	var c Course
	found := rows.Next()
	if found {
		if err := scanCourse(&c, rows); err != nil {
			rows.Close()
			return nil, err
		}
	}
	err = rows.Err()
	rows.Close()
	if !found {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	d.courseCounts(&c)
	return &c, nil
}

type CourseFilter struct {
	Search     string
	SubjectID  string
	TeacherID  string
	Department string
	Stage      string
	Status     string
	Page       int
	Limit      int
}

func (d *DB) ListCourses(f CourseFilter) ([]Course, int, error) {
	var conds []string
	var args []any
	if f.Search != "" {
		conds = append(conds, "c.title LIKE ?")
		args = append(args, "%"+f.Search+"%")
	}
	if f.SubjectID != "" {
		conds = append(conds, "c.subject_id = ?")
		args = append(args, f.SubjectID)
	}
	if f.TeacherID != "" {
		conds = append(conds, "c.teacher_id = ?")
		args = append(args, f.TeacherID)
	}
	if f.Department != "" {
		conds = append(conds, "c.department = ?")
		args = append(args, f.Department)
	}
	if f.Stage != "" {
		conds = append(conds, "c.stage = ?")
		args = append(args, f.Stage)
	}
	if f.Status != "" {
		conds = append(conds, "c.status = ?")
		args = append(args, f.Status)
	}
	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}
	var total int
	if err := d.SQL.QueryRow("SELECT COUNT(*) FROM courses c "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	limit := f.Limit
	if limit <= 0 || limit > 200 {
		limit = 20
	}
	page := f.Page
	if page < 1 {
		page = 1
	}
	rows, err := d.SQL.Query(courseSelect+" "+where+" ORDER BY c.created_at DESC LIMIT ? OFFSET ?",
		append(args, limit, (page-1)*limit)...)
	if err != nil {
		return nil, 0, err
	}
	var out []Course
	for rows.Next() {
		var c Course
		if err := scanCourse(&c, rows); err != nil {
			rows.Close()
			return nil, 0, err
		}
		out = append(out, c)
	}
	qerr := rows.Err()
	rows.Close()
	if qerr != nil {
		return nil, 0, qerr
	}
	for i := range out {
		d.courseCounts(&out[i])
	}
	return out, total, nil
}

func (d *DB) UpdateCourse(id string, fields map[string]any) (*Course, error) {
	allowed := map[string]bool{"title": true, "description": true, "thumbnail_url": true, "subject_id": true,
		"department": true, "department_id": true, "stage": true, "teacher_id": true, "price": true, "status": true}
	var sets []string
	var args []any
	for k, v := range fields {
		if allowed[k] {
			sets = append(sets, k+" = ?")
			args = append(args, v)
		}
	}
	if len(sets) == 0 {
		return d.GetCourse(id)
	}
	if st, ok := fields["status"]; ok && st == "published" {
		sets = append(sets, "is_published = 1", "published_at = ?")
		args = append(args, NowUTC())
	}
	sets = append(sets, "updated_at = ?")
	args = append(args, NowUTC(), id)
	_, err := d.SQL.Exec("UPDATE courses SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...)
	if err != nil {
		return nil, err
	}
	return d.GetCourse(id)
}

func (d *DB) SetCourseStatus(id, status string) (*Course, error) {
	published := 0
	publishedAt := ""
	if status == "published" {
		published, publishedAt = 1, NowUTC()
	}
	_, err := d.SQL.Exec("UPDATE courses SET status=?, is_published=?, published_at=?, updated_at=? WHERE id=?",
		status, published, publishedAt, NowUTC(), id)
	if err != nil {
		return nil, err
	}
	return d.GetCourse(id)
}

func (d *DB) DeleteCourse(id string) error {
	_, err := d.SQL.Exec("DELETE FROM courses WHERE id=?", id)
	if err != nil {
		return err
	}
	_, _ = d.SQL.Exec("DELETE FROM sections WHERE course_id=?", id)
	_, _ = d.SQL.Exec("DELETE FROM lessons WHERE course_id=?", id)
	_, _ = d.SQL.Exec("DELETE FROM enrollments WHERE course_id=?", id)
	return nil
}

func (d *DB) AssignCourseTeacher(courseID, teacherID string) (*Course, error) {
	_, err := d.SQL.Exec("UPDATE courses SET teacher_id=?, updated_at=? WHERE id=?", teacherID, NowUTC(), courseID)
	if err != nil {
		return nil, err
	}
	return d.GetCourse(courseID)
}

// ---------------------------------------------------------------------------
// Sections
// ---------------------------------------------------------------------------

const sectionCols = "id, course_id, title, description, sort_order, created_at, updated_at"

func scanSection(s *Section, rows *sql.Rows) error {
	err := rows.Scan(&s.ID, &s.CourseID, &s.Title, &s.Description, &s.Order, &s.CreatedAt, &s.UpdatedAt)
	if err == nil {
		_ = 0
	}
	return err
}

func (d *DB) SectionWithCount(id string) (*Section, error) {
	var s Section
	err := d.SQL.QueryRow("SELECT "+sectionCols+" FROM sections WHERE id=?", id).
		Scan(&s.ID, &s.CourseID, &s.Title, &s.Description, &s.Order, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM lessons WHERE section_id=?", id).Scan(&s.LessonsCount)
	return &s, nil
}

func (d *DB) ListSections(courseID string) ([]Section, error) {
	rows, err := d.SQL.Query("SELECT "+sectionCols+" FROM sections WHERE course_id=? ORDER BY sort_order, created_at", courseID)
	if err != nil {
		return nil, err
	}
	var out []Section
	for rows.Next() {
		var s Section
		if err := rows.Scan(&s.ID, &s.CourseID, &s.Title, &s.Description, &s.Order, &s.CreatedAt, &s.UpdatedAt); err != nil {
			rows.Close()
			return nil, err
		}
		out = append(out, s)
	}
	qerr := rows.Err()
	rows.Close()
	if qerr != nil {
		return nil, qerr
	}
	for i := range out {
		_ = d.SQL.QueryRow("SELECT COUNT(*) FROM lessons WHERE section_id=?", out[i].ID).Scan(&out[i].LessonsCount)
	}
	return out, nil
}

func (d *DB) CreateSection(courseID, title, description string, order int) (*Section, error) {
	if order == 0 {
		_ = d.SQL.QueryRow("SELECT COALESCE(MAX(sort_order),-1)+1 FROM sections WHERE course_id=?", courseID).Scan(&order)
	}
	id := NewID()
	now := NowUTC()
	_, err := d.SQL.Exec(`INSERT INTO sections(id,course_id,title,description,sort_order,created_at,updated_at)
		VALUES(?,?,?,?,?,?,?)`, id, courseID, title, description, order, now, now)
	if err != nil {
		return nil, err
	}
	return d.SectionWithCount(id)
}

func (d *DB) UpdateSection(id string, fields map[string]any) (*Section, error) {
	allowed := map[string]bool{"title": true, "description": true, "sort_order": true, "order": true}
	var sets []string
	var args []any
	for k, v := range fields {
		if k == "order" {
			k = "sort_order"
		}
		if allowed[k] {
			sets = append(sets, k+" = ?")
			args = append(args, v)
		}
	}
	if len(sets) > 0 {
		sets = append(sets, "updated_at = ?")
		args = append(args, NowUTC(), id)
		_, err := d.SQL.Exec("UPDATE sections SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...)
		if err != nil {
			return nil, err
		}
	}
	return d.SectionWithCount(id)
}

func (d *DB) DeleteSection(id string) error {
	_, _ = d.SQL.Exec("UPDATE lessons SET section_id='' WHERE section_id=?", id)
	_, err := d.SQL.Exec("DELETE FROM sections WHERE id=?", id)
	return err
}

func (d *DB) ReorderSections(courseID string, ids []string) ([]Section, error) {
	for i, id := range ids {
		_, _ = d.SQL.Exec("UPDATE sections SET sort_order=?, updated_at=? WHERE id=? AND course_id=?", i, NowUTC(), id, courseID)
	}
	return d.ListSections(courseID)
}

// ---------------------------------------------------------------------------
// Lessons
// ---------------------------------------------------------------------------

const lessonCols = "id, course_id, section_id, title, description, thumbnail, file_path, youtube_video_id, youtube_url, duration_seconds, sort_order, status, is_free, is_free_preview, views, created_by, published_at, created_at, updated_at"

func scanLesson(l *Lesson, rows *sql.Rows) error {
	var free, preview int
	var filePath string
	err := rows.Scan(&l.ID, &l.CourseID, &l.SectionID, &l.Title, &l.Description, &l.Thumbnail,
		&filePath, &l.YoutubeVideoID, &l.YoutubeURL, &l.Duration, &l.Order, &l.Status, &free, &preview,
		&l.Views, &l.CreatedBy, &l.PublishedAt, &l.CreatedAt, &l.UpdatedAt)
	_ = preview
	l.FilePath = filePath
	l.IsFree = free != 0
	return err
}

func (d *DB) GetLesson(id string) (*Lesson, error) {
	rows, err := d.SQL.Query("SELECT "+lessonCols+" FROM lessons WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	var l Lesson
	if err := scanLesson(&l, rows); err != nil {
		return nil, err
	}
	return &l, rows.Err()
}

func (d *DB) listLessons(where string, args ...any) ([]Lesson, error) {
	rows, err := d.SQL.Query("SELECT "+lessonCols+" FROM lessons "+where+" ORDER BY sort_order, created_at", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Lesson
	for rows.Next() {
		var l Lesson
		if err := scanLesson(&l, rows); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (d *DB) ListLessonsByCourse(courseID string) ([]Lesson, error) {
	return d.listLessons("WHERE course_id=?", courseID)
}

func (d *DB) ListLessonsBySection(sectionID string) ([]Lesson, error) {
	return d.listLessons("WHERE section_id=?", sectionID)
}

func (d *DB) CreateLesson(courseID, sectionID, title, description, thumbnail, youtubeVideoID, youtubeURL string, duration, order int, status string, isFree bool, createdBy string) (*Lesson, error) {
	if status == "" {
		status = "draft"
	}
	if order == 0 {
		_ = d.SQL.QueryRow("SELECT COALESCE(MAX(sort_order),-1)+1 FROM lessons WHERE course_id=? AND section_id=?", courseID, sectionID).Scan(&order)
	}
	id := NewID()
	now := NowUTC()
	free := 0
	if isFree {
		free = 1
	}
	_, err := d.SQL.Exec(`INSERT INTO lessons(id,course_id,section_id,title,description,thumbnail,youtube_video_id,youtube_url,duration_seconds,sort_order,status,is_free,created_by,created_at,updated_at)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, id, courseID, sectionID, title, description, thumbnail, youtubeVideoID, youtubeURL, duration, order, status, free, createdBy, now, now)
	if err != nil {
		return nil, err
	}
	return d.GetLesson(id)
}

func (d *DB) UpdateLesson(id string, fields map[string]any) (*Lesson, error) {
	allowed := map[string]bool{"title": true, "description": true, "thumbnail": true, "section_id": true,
		"youtube_video_id": true, "youtube_url": true, "duration_seconds": true, "duration": true,
		"sort_order": true, "order": true, "status": true, "is_free": true, "file_path": true}
	var sets []string
	var args []any
	for k, v := range fields {
		switch k {
		case "duration":
			k = "duration_seconds"
		case "order":
			k = "sort_order"
		case "is_free":
			if b, ok := v.(bool); ok {
				if b {
					v = 1
				} else {
					v = 0
				}
			}
		}
		if allowed[k] {
			sets = append(sets, k+" = ?")
			args = append(args, v)
		}
	}
	if len(sets) == 0 {
		return d.GetLesson(id)
	}
	if st, ok := fields["status"]; ok && st == "published" {
		sets = append(sets, "published_at = ?")
		args = append(args, NowUTC())
	}
	sets = append(sets, "updated_at = ?")
	args = append(args, NowUTC(), id)
	_, err := d.SQL.Exec("UPDATE lessons SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...)
	if err != nil {
		return nil, err
	}
	return d.GetLesson(id)
}

func (d *DB) SetLessonStatus(id, status string) (*Lesson, error) {
	publishedAt := ""
	if status == "published" {
		publishedAt = NowUTC()
		_, err := d.SQL.Exec("UPDATE lessons SET status=?, published_at=?, updated_at=? WHERE id=?", status, publishedAt, NowUTC(), id)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := d.SQL.Exec("UPDATE lessons SET status=?, updated_at=? WHERE id=?", status, NowUTC(), id)
		if err != nil {
			return nil, err
		}
	}
	return d.GetLesson(id)
}

func (d *DB) DeleteLesson(id string) error {
	_, _ = d.SQL.Exec("DELETE FROM attachments WHERE lesson_id=?", id)
	_, _ = d.SQL.Exec("DELETE FROM progress WHERE lesson_id=?", id)
	_, err := d.SQL.Exec("DELETE FROM lessons WHERE id=?", id)
	return err
}

func (d *DB) ReorderLessons(sectionID string, ids []string) ([]Lesson, error) {
	for i, id := range ids {
		_, _ = d.SQL.Exec("UPDATE lessons SET sort_order=?, updated_at=? WHERE id=? AND section_id=?", i, NowUTC(), id, sectionID)
	}
	return d.ListLessonsBySection(sectionID)
}

func (d *DB) SetLessonVideo(id, filePath string) (*Lesson, error) {
	_, err := d.SQL.Exec("UPDATE lessons SET file_path=?, updated_at=? WHERE id=?", filePath, NowUTC(), id)
	if err != nil {
		return nil, err
	}
	return d.GetLesson(id)
}

// ---------------------------------------------------------------------------
// Attachments
// ---------------------------------------------------------------------------

func (d *DB) ListAttachments(lessonID string) ([]Attachment, error) {
	rows, err := d.SQL.Query(`SELECT id, lesson_id, name, type, size_bytes, url, storage_key, created_by, created_at
		FROM attachments WHERE lesson_id=? ORDER BY created_at`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Attachment
	for rows.Next() {
		var a Attachment
		if err := rows.Scan(&a.ID, &a.LessonID, &a.Name, &a.Type, &a.Size, &a.URL, &a.StorageKey, &a.CreatedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (d *DB) CreateAttachment(lessonID, name, typ, url, storageKey, mimeType string, sizeBytes int64, createdBy string) (*Attachment, error) {
	if typ == "" {
		typ = "file"
	}
	a := &Attachment{ID: NewID(), LessonID: lessonID, Name: name, Type: typ, Size: sizeBytes,
		URL: url, StorageKey: storageKey, CreatedBy: createdBy, CreatedAt: NowUTC()}
	_, err := d.SQL.Exec(`INSERT INTO attachments(id,lesson_id,name,type,size_bytes,url,storage_key,mime_type,created_by,created_at)
		VALUES(?,?,?,?,?,?,?,?,?,?)`, a.ID, a.LessonID, a.Name, a.Type, a.Size, a.URL, a.StorageKey, mimeType, a.CreatedBy, a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (d *DB) DeleteAttachment(id string) error {
	_, err := d.SQL.Exec("DELETE FROM attachments WHERE id=?", id)
	return err
}

// ---------------------------------------------------------------------------
// Enrollments & subscriptions
// ---------------------------------------------------------------------------

func (d *DB) Enroll(studentID, courseID, status string) (*Enrollment, error) {
	if status == "" {
		status = "active"
	}
	now := NowUTC()
	e := &Enrollment{ID: NewID(), StudentID: studentID, CourseID: courseID, Status: status, CreatedAt: now}
	_, err := d.SQL.Exec(`INSERT INTO enrollments(id,student_id,course_id,status,created_at)
		VALUES(?,?,?,?,?) ON CONFLICT(student_id, course_id) DO UPDATE SET status=excluded.status`,
		e.ID, e.StudentID, e.CourseID, e.Status, now)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (d *DB) Unenroll(studentID, courseID string) error {
	_, err := d.SQL.Exec("DELETE FROM enrollments WHERE student_id=? AND course_id=?", studentID, courseID)
	return err
}

func (d *DB) ListEnrollmentsByStudent(studentID string) ([]Enrollment, error) {
	rows, err := d.SQL.Query(`SELECT e.id, e.student_id, e.course_id, COALESCE(c.title,''), e.status, e.expires_at, e.created_at
		FROM enrollments e LEFT JOIN courses c ON c.id=e.course_id WHERE e.student_id=? ORDER BY e.created_at DESC`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(&e.ID, &e.StudentID, &e.CourseID, &e.CourseTitle, &e.Status, &e.ExpiresAt, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (d *DB) ListEnrollmentsByCourse(courseID string) ([]Enrollment, error) {
	rows, err := d.SQL.Query(`SELECT id, student_id, course_id, '', status, expires_at, created_at
		FROM enrollments WHERE course_id=? ORDER BY created_at DESC`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(&e.ID, &e.StudentID, &e.CourseID, &e.CourseTitle, &e.Status, &e.ExpiresAt, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (d *DB) ListSubscriptionsByStudent(studentID string) ([]Enrollment, error) {
	rows, err := d.SQL.Query(`SELECT id, student_id, course_id, '', status, expires_at, created_at
		FROM subscriptions WHERE student_id=? ORDER BY created_at DESC`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(&e.ID, &e.StudentID, &e.CourseID, &e.CourseTitle, &e.Status, &e.ExpiresAt, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// HasAccess reports whether the student may watch the course (free or enrolled/subscribed).
func (d *DB) HasAccess(studentID, courseID string, free bool) bool {
	if free {
		return true
	}
	if courseID == "" || studentID == "" {
		return false
	}
	now := NowUTC()
	var n int
	_ = d.SQL.QueryRow(`SELECT COUNT(*) FROM enrollments WHERE student_id=? AND course_id=? AND status='active'
		AND (expires_at='' OR expires_at IS NULL OR expires_at > ?)`, studentID, courseID, now).Scan(&n)
	if n > 0 {
		return true
	}
	_ = d.SQL.QueryRow(`SELECT COUNT(*) FROM subscriptions WHERE student_id=? AND course_id=? AND status='active'
		AND (expires_at='' OR expires_at IS NULL OR expires_at > ?)`, studentID, courseID, now).Scan(&n)
	return n > 0
}

// ---------------------------------------------------------------------------
// Progress
// ---------------------------------------------------------------------------

type LessonProgress struct {
	LessonID       string `json:"lesson_id"`
	StudentID      string `json:"student_id"`
	WatchedSeconds int    `json:"watched_seconds"`
	Completed      bool   `json:"completed"`
	UpdatedAt      string `json:"updated_at"`
}

func (d *DB) ListProgressByStudent(studentID string) ([]LessonProgress, error) {
	rows, err := d.SQL.Query(`SELECT lesson_id, student_id, watched_seconds, completed, updated_at FROM progress WHERE student_id=?`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []LessonProgress
	for rows.Next() {
		var p LessonProgress
		var completed int
		if err := rows.Scan(&p.LessonID, &p.StudentID, &p.WatchedSeconds, &completed, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Completed = completed != 0
		out = append(out, p)
	}
	return out, rows.Err()
}

type CourseProgress struct {
	CourseID         string `json:"course_id"`
	TotalLessons     int    `json:"total_lessons"`
	CompletedLessons int    `json:"completed_lessons"`
	Percent          int    `json:"percent"`
}

func (d *DB) CourseProgress(studentID, courseID string) (*CourseProgress, error) {
	cp := &CourseProgress{CourseID: courseID}
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM lessons WHERE course_id=?", courseID).Scan(&cp.TotalLessons)
	_ = d.SQL.QueryRow(`SELECT COUNT(*) FROM progress p JOIN lessons l ON l.id=p.lesson_id
		WHERE p.student_id=? AND l.course_id=? AND p.completed=1`, studentID, courseID).Scan(&cp.CompletedLessons)
	if cp.TotalLessons > 0 {
		cp.Percent = cp.CompletedLessons * 100 / cp.TotalLessons
	}
	return cp, nil
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

type Overview struct {
	TotalStudents       int `json:"total_students"`
	ActiveStudents      int `json:"active_students"`
	TotalCourses        int `json:"total_courses"`
	TotalSubjects       int `json:"total_subjects"`
	TotalLessons        int `json:"total_lessons"`
	TotalViews          int `json:"total_views"`
	ActiveSubscriptions int `json:"active_subscriptions"`
	UsedSubscriptionCodes int `json:"used_subscription_codes"`
}

func (d *DB) Overview() (*Overview, error) {
	o := &Overview{}
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM users WHERE role='student'").Scan(&o.TotalStudents)
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM users WHERE role='student' AND is_active=1").Scan(&o.ActiveStudents)
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM courses").Scan(&o.TotalCourses)
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM subjects").Scan(&o.TotalSubjects)
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM lessons").Scan(&o.TotalLessons)
	_ = d.SQL.QueryRow("SELECT COALESCE(SUM(views),0) FROM lessons").Scan(&o.TotalViews)
	_ = d.SQL.QueryRow("SELECT COUNT(*) FROM subscriptions WHERE status='active'").Scan(&o.ActiveSubscriptions)
	return o, nil
}

func (d *DB) RecentCourses(limit int) ([]Course, error) {
	if limit <= 0 {
		limit = 5
	}
	rows, err := d.SQL.Query(courseSelect+" ORDER BY c.created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	var out []Course
	for rows.Next() {
		var c Course
		if err := scanCourse(&c, rows); err != nil {
			rows.Close()
			return nil, err
		}
		out = append(out, c)
	}
	qerr := rows.Err()
	rows.Close()
	if qerr != nil {
		return nil, qerr
	}
	for i := range out {
		d.courseCounts(&out[i])
	}
	return out, nil
}

func (d *DB) RecentStudents(limit int) ([]User, error) {
	if limit <= 0 {
		limit = 5
	}
	rows, err := d.SQL.Query("SELECT "+userCols+" FROM users WHERE role='student' ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		var u User
		if err := scanUserRows(&u, rows); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// ---------------------------------------------------------------------------
// Audit log
// ---------------------------------------------------------------------------

func (d *DB) Audit(actorID, action, entity, entityID, meta string) {
	_, _ = d.SQL.Exec(`INSERT INTO audit_logs(id,actor_id,action,entity,entity_id,meta,created_at)
		VALUES(?,?,?,?,?,?,?)`, NewID(), actorID, action, entity, entityID, meta, NowUTC())
}

func (d *DB) RecentActivity(limit int) ([]ActivityItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := d.SQL.Query(`SELECT id, action || ' ' || entity, COALESCE(meta,''), created_at FROM audit_logs ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ActivityItem
	for rows.Next() {
		var a ActivityItem
		if err := rows.Scan(&a.ID, &a.Type, &a.Message, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
