package video

import (
	"context"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/db"
)

// SQLiteCatalog implements Catalog against the local SQLite database.
// It is the default persistent catalog: lessons carry the playable video
// source (file_path), and access is decided from enrollments/subscriptions.
type SQLiteCatalog struct {
	db *db.DB
}

// NewSQLiteCatalog opens (creating if needed) the SQLite database at path.
func NewSQLiteCatalog(path string) (*SQLiteCatalog, error) {
	d, err := db.Open(path)
	if err != nil {
		return nil, err
	}
	return &SQLiteCatalog{db: d}, nil
}

// DB exposes the underlying store for admin handlers.
func (c *SQLiteCatalog) DB() *db.DB { return c.db }

func (c *SQLiteCatalog) Close() error { return c.db.Close() }

func (c *SQLiteCatalog) Find(ctx context.Context, id string) (Item, error) {
	// Lessons are the primary playable entity.
	if l, err := c.db.GetLesson(id); err == nil {
		return Item{
			ID:         l.ID,
			Title:      l.Title,
			CourseID:   l.CourseID,
			SourcePath: l.FilePath,
			YouTubeURL: l.YoutubeURL,
			YouTubeID:  l.YoutubeVideoID,
			MimeType:   "video/mp4",
			Status:     l.Status,
			Published:  l.Status == "published",
			Free:       l.IsFree,
		}, nil
	}
	// Fall back to the videos table (upload pipeline).
	var item Item
	err := c.db.SQL.QueryRowContext(ctx, `SELECT id, title, COALESCE(course_id,''),
		COALESCE(file_path, temp_path, ''), COALESCE(youtube_url,''), COALESCE(youtube_video_id,''),
		COALESCE(mime_type,'video/mp4'), COALESCE(status,'pending'),
		COALESCE(is_published,0), 0
		FROM videos WHERE id = ? LIMIT 1`, id).Scan(
		&item.ID, &item.Title, &item.CourseID, &item.SourcePath, &item.YouTubeURL,
		&item.YouTubeID, &item.MimeType, &item.Status, &item.Published, &item.Free)
	if err != nil {
		return Item{}, ErrNotFound
	}
	return item, nil
}

func (c *SQLiteCatalog) CanAccess(ctx context.Context, userID string, item Item) (bool, error) {
	// Resolve the caller's internal user id: identity IDs are telegram ids.
	studentID := userID
	var dbID string
	_ = c.db.SQL.QueryRowContext(ctx, `SELECT id FROM users WHERE id = ? OR CAST(telegram_id AS TEXT) = ? LIMIT 1`,
		userID, userID).Scan(&dbID)
	if dbID != "" {
		studentID = dbID
	}
	return c.db.HasAccess(studentID, item.CourseID, item.Free), nil
}
