package video

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/netutil"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// DialFunc uses the shared IPv4-only dialer from netutil for every
// connection this catalog opens (see internal/netutil).
var dialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
	return netutil.Ipv4Dialer().DialContext(ctx, network, addr)
}

type Item struct {
	ID         string
	Title      string
	CourseID   string
	SourcePath string
	YouTubeURL string
	YouTubeID  string
	MimeType   string
	Status     string
	Published  bool
	Free       bool
}

type Catalog interface {
	Close() error
	Find(context.Context, string) (Item, error)
	CanAccess(context.Context, string, Item) (bool, error)
}

var ErrNotFound = errors.New("video not found")

// SQLCatalog works with the existing Supabase PostgreSQL database. The videos
// migration creates the table; lesson fallback keeps older records playable
// while they are being migrated.
type SQLCatalog struct {
	db   *sql.DB
	pool *pgxpool.Pool
}

// openDBFromPool wraps a *pgxpool.Pool as an *sql.DB so the rest of the code
// keeps using the standard library query API while the pool honors the
// custom IPv4-preferring dialer.
func openDBFromPool(pool *pgxpool.Pool) *sql.DB {
	return stdlib.OpenDBFromPool(pool)
}

func NewSQL(databaseURL string, maxOpen, maxIdle int, maxLifetime time.Duration) (*SQLCatalog, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, errors.New("DATABASE_URL is not configured")
	}
	// pgx stdlib's sql.Open does not expose the dial function, so parse a
	// pool config with an IPv4-preferring DialFunc and reuse the pool (which
	// implements the interfaces *sql.DB needs) for queries.
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	// Supabase direct database hosts (db.<ref>.supabase.co) are IPv6-only;
	// hosts without IPv6 egress (e.g. Render Oregon) can never reach them.
	// Supavisor pooler hosts are always IPv4, so rewrite direct Supabase
	// connection strings to the session-mode pooler automatically.
	if isSupabaseDirectHost(cfg.ConnConfig.Host) {
		if err := rewriteToSupavisor(cfg); err != nil {
			return nil, err
		}
	} else if h := cfg.ConnConfig.Host; h != "" {
		// For every other host, resolve to an IPv4 address first so the
		// connection string carries a numeric IP and internal resolvers
		// never return an AAAA record.
		if v4, err := netutil.HostPortIPv4(context.Background(), net.JoinHostPort(h, "5432")); err == nil {
			if host, _, splitErr := net.SplitHostPort(v4); splitErr == nil {
				cfg.ConnConfig.Host = host
			}
		}
	}
	cfg.ConnConfig.DialFunc = dialFunc
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	db := openDBFromPool(pool)
	if maxOpen > 0 {
		db.SetMaxOpenConns(maxOpen)
	}
	if maxIdle > 0 {
		db.SetMaxIdleConns(maxIdle)
	}
	if maxLifetime > 0 {
		db.SetConnMaxLifetime(maxLifetime)
	}
	if err := db.PingContext(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &SQLCatalog{db: db, pool: pool}, nil
}
func (c *SQLCatalog) Close() error { c.pool.Close(); return c.db.Close() }

// isSupabaseDirectHost reports whether host is a direct Supabase database
// host (db.<ref>.supabase.co) which is IPv6-only.
func isSupabaseDirectHost(host string) bool {
	return strings.HasSuffix(strings.TrimSpace(host), ".supabase.co")
}

// rewriteToSupavisor converts a direct Supabase connection config to the
// session-mode Supavisor pooler, which is always reachable over IPv4. The
// pooler requires the username to be "postgres.<project-ref>".
func rewriteToSupavisor(cfg *pgxpool.Config) error {
	ref := extractProjectRef(cfg.ConnConfig.Host)
	if ref == "" {
		return errors.New("unable to extract Supabase project ref from " + cfg.ConnConfig.Host)
	}
	cfg.ConnConfig.User = "postgres." + ref
	// The Supavisor host follows the pattern aws-0-<region>.pooler.supabase.com.
	// If the project lives outside us-east-1, override SUPABASE_POOLER_HOST
	// (e.g. aws-0-ap-southeast-1.pooler.supabase.com) in the environment.
	if h := os.Getenv("SUPABASE_POOLER_HOST"); h != "" {
		cfg.ConnConfig.Host = h
	} else {
		cfg.ConnConfig.Host = "aws-0-us-east-1.pooler.supabase.com"
	}
	cfg.ConnConfig.Port = 5432
	return nil
}

// extractProjectRef extracts <ref> from "db.<ref>.supabase.co".
func extractProjectRef(host string) string {
	withoutDB := strings.TrimPrefix(host, "db.")
	ref, _, ok := strings.Cut(withoutDB, ".supabase.co")
	if !ok || ref == "" {
		return ""
	}
	return ref
}

func (c *SQLCatalog) Find(ctx context.Context, id string) (Item, error) {
	var item Item
	err := c.db.QueryRowContext(ctx, `
		SELECT id::text, title, COALESCE(course_id, ''), COALESCE(file_path, temp_path, ''),
		       COALESCE(youtube_url, ''), COALESCE(youtube_video_id, ''),
		       COALESCE(mime_type, 'video/mp4'), COALESCE(status, 'pending'),
		       COALESCE(is_published, false), false
		FROM public.videos WHERE id::text = $1 LIMIT 1`, id).Scan(
		&item.ID, &item.Title, &item.CourseID, &item.SourcePath, &item.YouTubeURL,
		&item.YouTubeID, &item.MimeType, &item.Status, &item.Published, &item.Free,
	)
	if err == nil {
		return item, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Item{}, err
	}

	err = c.db.QueryRowContext(ctx, `
		SELECT id::text, title, COALESCE(course_id, ''), '',
		       COALESCE(youtube_url, ''), COALESCE(youtube_video_id, ''),
		       'video/mp4', COALESCE(status, 'published'),
		       (COALESCE(status, '') = 'published'),
		       (COALESCE(is_free, false) OR COALESCE(is_free_preview, false))
		FROM public.lessons WHERE id::text = $1 LIMIT 1`, id).Scan(
		&item.ID, &item.Title, &item.CourseID, &item.SourcePath, &item.YouTubeURL,
		&item.YouTubeID, &item.MimeType, &item.Status, &item.Published, &item.Free,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Item{}, ErrNotFound
	}
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func (c *SQLCatalog) CanAccess(ctx context.Context, userID string, item Item) (bool, error) {
	if item.Free {
		return true, nil
	}
	if item.CourseID == "" {
		return false, nil
	}
	var allowed bool
	err := c.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM public.enrollments e
			WHERE e.student_id = $1 AND e.course_id = $2 AND e.status = 'active'
			  AND (e.expires_at IS NULL OR e.expires_at > now())
		) OR EXISTS (
			SELECT 1 FROM public.subscriptions s
			WHERE s.student_id = $1 AND s.course_id = $2 AND s.status = 'active'
			  AND (s.expires_at IS NULL OR s.expires_at > now())
		)`, userID, item.CourseID).Scan(&allowed)
	return allowed, err
}

// MemoryCatalog is intentionally only for local development and tests. It
// never reads paths outside the configured upload directory.
type MemoryCatalog struct {
	items     map[string]Item
	uploadDir string
}

func NewMemory(uploadDir string) *MemoryCatalog {
	return &MemoryCatalog{items: map[string]Item{}, uploadDir: uploadDir}
}
func (c *MemoryCatalog) Close() error  { return nil }
func (c *MemoryCatalog) Put(item Item) { c.items[item.ID] = item }
func (c *MemoryCatalog) Find(_ context.Context, id string) (Item, error) {
	item, ok := c.items[id]
	if !ok {
		return Item{}, ErrNotFound
	}
	return item, nil
}
func (c *MemoryCatalog) CanAccess(_ context.Context, _ string, item Item) (bool, error) {
	return item.Free, nil
}
func SafePath(root, candidate string) (string, error) {
	if strings.TrimSpace(root) == "" || strings.TrimSpace(candidate) == "" {
		return "", errors.New("video source is empty")
	}
	if strings.Contains(candidate, "\x00") {
		return "", errors.New("invalid video path")
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	// The candidate is stored relative to the upload root (e.g.
	// "videos/<id>_name.mp4"), so join it with the root before resolving.
	joined := candidate
	if !filepath.IsAbs(candidate) {
		joined = filepath.Join(root, candidate)
	}
	candidateAbs, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	rootReal, err := filepath.EvalSymlinks(rootAbs)
	if err != nil {
		return "", err
	}
	candidateReal, err := filepath.EvalSymlinks(candidateAbs)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(rootReal, candidateReal)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", errors.New("video source is outside upload directory")
	}
	info, err := os.Stat(candidateReal)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", errors.New("video source is a directory")
	}
	return candidateReal, nil
}
