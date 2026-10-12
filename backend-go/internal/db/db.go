// Package db provides the SQLite-backed persistent store for Wiyana.
//
// It replaces the previous PostgreSQL/Supabase + Convex backends: a single
// local SQLite file holds users, courses, lessons, enrollments and the video
// catalog. The driver is modernc.org/sqlite (pure Go, no CGO).
package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite connection.
type DB struct {
	SQL *sql.DB
}

// Open creates (or opens) the SQLite file at path, creating parent
// directories as needed, and ensures the schema exists.
func Open(path string) (*DB, error) {
	if path == "" {
		path = "data/wiyana.db"
	}
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	sqlDB, err := sql.Open("sqlite", path+"?cache=shared")
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	d := &DB{SQL: sqlDB}
	if err := d.migrate(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return d, nil
}

// Close releases the database.
func (d *DB) Close() error { return d.SQL.Close() }

// NewID generates a random 16-byte hex identifier.
func NewID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// NowUTC returns the current UTC time formatted for storage.
func NowUTC() string { return time.Now().UTC().Format(time.RFC3339) }

const schema = `
CREATE TABLE IF NOT EXISTS users(
  id TEXT PRIMARY KEY,
  telegram_id INTEGER UNIQUE,
  username TEXT DEFAULT '',
  first_name TEXT DEFAULT '',
  last_name TEXT DEFAULT '',
  photo_url TEXT DEFAULT '',
  role TEXT DEFAULT 'student',
  is_active INTEGER DEFAULT 1,
  phone TEXT DEFAULT '',
  bio TEXT DEFAULT '',
  created_at TEXT,
  updated_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

CREATE TABLE IF NOT EXISTS subjects(
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  code TEXT UNIQUE,
  description TEXT DEFAULT '',
  department TEXT DEFAULT '',
  department_id TEXT DEFAULT '',
  stage TEXT DEFAULT '',
  status TEXT DEFAULT 'active',
  created_by TEXT DEFAULT '',
  created_at TEXT,
  updated_at TEXT
);

CREATE TABLE IF NOT EXISTS courses(
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT DEFAULT '',
  thumbnail_url TEXT DEFAULT '',
  subject_id TEXT DEFAULT '',
  department TEXT DEFAULT '',
  department_id TEXT DEFAULT '',
  stage TEXT DEFAULT '',
  teacher_id TEXT DEFAULT '',
  price REAL DEFAULT 0,
  status TEXT DEFAULT 'draft',
  is_published INTEGER DEFAULT 0,
  published_at TEXT DEFAULT '',
  created_by TEXT DEFAULT '',
  created_at TEXT,
  updated_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_courses_subject ON courses(subject_id);
CREATE INDEX IF NOT EXISTS idx_courses_teacher ON courses(teacher_id);
CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);

CREATE TABLE IF NOT EXISTS sections(
  id TEXT PRIMARY KEY,
  course_id TEXT NOT NULL,
  title TEXT NOT NULL,
  description TEXT DEFAULT '',
  sort_order INTEGER DEFAULT 0,
  created_at TEXT,
  updated_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_sections_course ON sections(course_id);

CREATE TABLE IF NOT EXISTS lessons(
  id TEXT PRIMARY KEY,
  course_id TEXT NOT NULL,
  section_id TEXT DEFAULT '',
  title TEXT NOT NULL,
  description TEXT DEFAULT '',
  thumbnail TEXT DEFAULT '',
  file_path TEXT DEFAULT '',
  youtube_video_id TEXT DEFAULT '',
  youtube_url TEXT DEFAULT '',
  duration_seconds INTEGER DEFAULT 0,
  sort_order INTEGER DEFAULT 0,
  status TEXT DEFAULT 'draft',
  is_free INTEGER DEFAULT 0,
  is_free_preview INTEGER DEFAULT 0,
  views INTEGER DEFAULT 0,
  created_by TEXT DEFAULT '',
  published_at TEXT DEFAULT '',
  created_at TEXT,
  updated_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_lessons_course ON lessons(course_id);
CREATE INDEX IF NOT EXISTS idx_lessons_section ON lessons(section_id);

CREATE TABLE IF NOT EXISTS attachments(
  id TEXT PRIMARY KEY,
  lesson_id TEXT NOT NULL,
  name TEXT NOT NULL,
  type TEXT DEFAULT 'file',
  size_bytes INTEGER DEFAULT 0,
  url TEXT DEFAULT '',
  storage_key TEXT DEFAULT '',
  mime_type TEXT DEFAULT '',
  created_by TEXT DEFAULT '',
  created_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_attachments_lesson ON attachments(lesson_id);

CREATE TABLE IF NOT EXISTS enrollments(
  id TEXT PRIMARY KEY,
  student_id TEXT NOT NULL,
  course_id TEXT NOT NULL,
  status TEXT DEFAULT 'active',
  expires_at TEXT DEFAULT '',
  created_at TEXT,
  UNIQUE(student_id, course_id)
);
CREATE INDEX IF NOT EXISTS idx_enrollments_student ON enrollments(student_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_course ON enrollments(course_id);

CREATE TABLE IF NOT EXISTS subscriptions(
  id TEXT PRIMARY KEY,
  student_id TEXT NOT NULL,
  course_id TEXT NOT NULL,
  status TEXT DEFAULT 'active',
  expires_at TEXT DEFAULT '',
  created_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_subscriptions_student ON subscriptions(student_id);

CREATE TABLE IF NOT EXISTS progress(
  id TEXT PRIMARY KEY,
  student_id TEXT NOT NULL,
  lesson_id TEXT NOT NULL,
  watched_seconds INTEGER DEFAULT 0,
  completed INTEGER DEFAULT 0,
  updated_at TEXT,
  UNIQUE(student_id, lesson_id)
);

CREATE TABLE IF NOT EXISTS videos(
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT DEFAULT '',
  file_path TEXT DEFAULT '',
  temp_path TEXT DEFAULT '',
  file_size INTEGER DEFAULT 0,
  mime_type TEXT DEFAULT 'video/mp4',
  status TEXT DEFAULT 'pending',
  uploader_id TEXT DEFAULT '',
  course_id TEXT DEFAULT '',
  youtube_video_id TEXT DEFAULT '',
  youtube_url TEXT DEFAULT '',
  is_published INTEGER DEFAULT 0,
  created_at TEXT,
  updated_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_videos_course ON videos(course_id);
CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);

CREATE TABLE IF NOT EXISTS audit_logs(
  id TEXT PRIMARY KEY,
  actor_id TEXT DEFAULT '',
  action TEXT NOT NULL,
  entity TEXT DEFAULT '',
  entity_id TEXT DEFAULT '',
  meta TEXT DEFAULT '',
  created_at TEXT
);
`

func (d *DB) migrate() error {
	_, err := d.SQL.Exec(schema)
	return err
}
