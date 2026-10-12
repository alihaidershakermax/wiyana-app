// SQLiteStore adapts the SQLite db.DB to the ConvexStore interface the
// Telegram bot expects, so the bot works without a Convex backend.
package bot

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/db"
)

// SQLiteStore implements ConvexStore on top of SQLite.
type SQLiteStore struct {
	db *db.DB

	mu         sync.Mutex
	shareLinks map[string]string // hash -> file_id (in-memory share links)
}

// NewSQLiteStore returns a ConvexStore backed by the given SQLite database.
func NewSQLiteStore(d *db.DB) *SQLiteStore {
	return &SQLiteStore{db: d, shareLinks: make(map[string]string)}
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func strArg(args map[string]any, key string) string {
	if args == nil {
		return ""
	}
	switch v := args[key].(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	default:
		return fmt.Sprint(v)
	}
}

func userToDoc(u *db.User) map[string]any {
	return map[string]any{
		"_id":         u.ID,
		"telegram_id": strconv.FormatInt(u.TelegramID, 10),
		"first_name":  u.FirstName,
		"is_admin":    u.Role == "admin" || u.Role == "owner",
		"is_active":   u.IsActive,
	}
}

func (s *SQLiteStore) getUserByTelegram(ctx context.Context, telegramID string) (*db.User, error) {
	tid, err := strconv.ParseInt(strings.TrimSpace(telegramID), 10, 64)
	if err != nil {
		return nil, nil
	}
	var u db.User
	row := s.db.SQL.QueryRowContext(ctx,
		`SELECT id, telegram_id, username, first_name, last_name, photo_url, role, is_active, phone, bio, created_at, updated_at
		 FROM users WHERE telegram_id = ? LIMIT 1`, tid)
	if err := scanUserRow(&u, row); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func scanUserRow(u *db.User, row *sql.Row) error {
	var active int
	err := row.Scan(&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.PhotoURL, &u.Role, &active, &u.Phone, &u.Bio, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}
	u.IsActive = active != 0
	return nil
}

// CallQuery implements ConvexStore.
func (s *SQLiteStore) CallQuery(ctx context.Context, path string, args map[string]any) (json.RawMessage, error) {
	switch path {
	case "users:getUser":
		u, err := s.getUserByTelegram(ctx, strArg(args, "telegramId"))
		if err != nil || u == nil {
			return nil, err
		}
		return mustJSON(userToDoc(u)), nil

	case "users:listUsers":
		users, _, err := s.db.ListUsers(db.UserFilter{})
		if err != nil {
			return nil, err
		}
		docs := make([]map[string]any, 0, len(users))
		for i := range users {
			docs = append(docs, userToDoc(&users[i]))
		}
		return mustJSON(docs), nil

	case "users:getShareLink":
		hash := strArg(args, "hash")
		s.mu.Lock()
		fileID := s.shareLinks[hash]
		s.mu.Unlock()
		if fileID == "" {
			return nil, nil
		}
		return mustJSON(map[string]any{"file_id": fileID}), nil

	case "users:searchVideos":
		q := strArg(args, "q")
		return s.searchLessons(ctx, q)

	case "users:listVideosBySubject":
		subject := strArg(args, "subject")
		return s.lessonsBySubject(ctx, subject)

	case "videos:listVideos":
		return s.allLessons(ctx)

	case "videos:getVideo":
		id := strArg(args, "id")
		l, err := s.db.GetLesson(id)
		if err != nil || l == nil {
			return nil, err
		}
		return mustJSON(lessonToVideoDoc(l)), nil

	case "users:listCategories":
		// Categories = distinct departments/stages from subjects.
		rows, err := s.db.SQL.QueryContext(ctx,
			`SELECT DISTINCT COALESCE(NULLIF(department,''), 'عام') FROM subjects ORDER BY 1`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var cats []map[string]any
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return nil, err
			}
			cats = append(cats, map[string]any{"name": name})
		}
		if cats == nil {
			cats = []map[string]any{}
		}
		return mustJSON(cats), nil

	case "users:listSubjects":
		category := strArg(args, "category")
		var subjects []db.Subject
		var err error
		if category == "" {
			subjects, _, err = s.db.ListSubjects(db.SubjectFilter{})
		} else {
			// Filter by department in Go (SubjectFilter may not cover it).
			all, _, err2 := s.db.ListSubjects(db.SubjectFilter{})
			if err2 != nil {
				return nil, err2
			}
			for _, sub := range all {
				dep := sub.Department
				if dep == "" {
					dep = "عام"
				}
				if dep == category {
					subjects = append(subjects, sub)
				}
			}
		}
		if err != nil {
			return nil, err
		}
		names := make([]map[string]any, 0, len(subjects))
		for _, sub := range subjects {
			names = append(names, map[string]any{"name": sub.Name})
		}
		return mustJSON(names), nil

	default:
		return nil, fmt.Errorf("unsupported query: %s", path)
	}
}

// CallMutation implements ConvexStore.
func (s *SQLiteStore) CallMutation(ctx context.Context, path string, args map[string]any) (json.RawMessage, error) {
	switch path {
	case "users:upsertUser":
		tid, _ := strconv.ParseInt(strArg(args, "telegram_id"), 10, 64)
		u, err := s.db.UpsertByTelegram(tid,
			strArg(args, "username"),
			strArg(args, "first_name"),
			strArg(args, "last_name"),
			"", "student")
		if err != nil {
			return nil, err
		}
		return mustJSON(userToDoc(u)), nil

	case "users:banUser", "users:unbanUser":
		u, err := s.getUserByTelegram(ctx, strArg(args, "telegramId"))
		if err != nil {
			return nil, err
		}
		if u == nil {
			return nil, fmt.Errorf("user not found")
		}
		if err := s.db.SetUserActive(u.ID, path == "users:unbanUser"); err != nil {
			return nil, err
		}
		return mustJSON(map[string]any{"ok": true}), nil

	case "users:createShareLink":
		fileID := strArg(args, "file_id")
		var b [8]byte
		if _, err := rand.Read(b[:]); err != nil {
			return nil, err
		}
		hash := hex.EncodeToString(b[:])
		s.mu.Lock()
		s.shareLinks[hash] = fileID
		s.mu.Unlock()
		return mustJSON(hash), nil

	case "users:incrementDownloads":
		// No download counter in SQLite schema; accept as no-op.
		return mustJSON(map[string]any{"ok": true}), nil

	default:
		return nil, fmt.Errorf("unsupported mutation: %s", path)
	}
}

func lessonToVideoDoc(l *db.Lesson) map[string]any {
	return map[string]any{
		"_id":             l.ID,
		"title":           l.Title,
		"category":        "",
		"subject":         "",
		"telegram_file_id": l.FilePath,
		"downloads":       0,
	}
}

func (s *SQLiteStore) searchLessons(ctx context.Context, q string) (json.RawMessage, error) {
	like := "%" + q + "%"
	rows, err := s.db.SQL.QueryContext(ctx,
		`SELECT id, course_id, section_id, title, description, thumbnail, file_path,
		        youtube_video_id, youtube_url, duration, "order", status, is_free,
		        views, created_by, published_at, created_at, updated_at
		 FROM lessons WHERE title LIKE ? OR description LIKE ? ORDER BY title LIMIT 50`, like, like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanLessonDocs(rows)
}

func (s *SQLiteStore) allLessons(ctx context.Context) (json.RawMessage, error) {
	rows, err := s.db.SQL.QueryContext(ctx,
		`SELECT id, course_id, section_id, title, description, thumbnail, file_path,
		        youtube_video_id, youtube_url, duration, "order", status, is_free,
		        views, created_by, published_at, created_at, updated_at
		 FROM lessons ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanLessonDocs(rows)
}

func (s *SQLiteStore) lessonsBySubject(ctx context.Context, subject string) (json.RawMessage, error) {
	like := "%" + subject + "%"
	rows, err := s.db.SQL.QueryContext(ctx,
		`SELECT l.id, l.course_id, l.section_id, l.title, l.description, l.thumbnail, l.file_path,
		        l.youtube_video_id, l.youtube_url, l.duration, l."order", l.status, l.is_free,
		        l.views, l.created_by, l.published_at, l.created_at, l.updated_at
		 FROM lessons l
		 JOIN courses c ON c.id = l.course_id
		 JOIN subjects sub ON sub.id = c.subject_id
		 WHERE sub.name LIKE ? ORDER BY l."order" LIMIT 100`, like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanLessonDocs(rows)
}

func (s *SQLiteStore) scanLessonDocs(rows *sql.Rows) (json.RawMessage, error) {
	docs := []map[string]any{}
	for rows.Next() {
		var l db.Lesson
		var isFree int
		if err := rows.Scan(&l.ID, &l.CourseID, &l.SectionID, &l.Title, &l.Description,
			&l.Thumbnail, &l.FilePath, &l.YoutubeVideoID, &l.YoutubeURL, &l.Duration,
			&l.Order, &l.Status, &isFree, &l.Views, &l.CreatedBy, &l.PublishedAt,
			&l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		l.IsFree = isFree != 0
		docs = append(docs, lessonToVideoDoc(&l))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mustJSON(docs), nil
}
