package video

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/netutil"
)

// ConvexCatalog implements the Catalog interface against a Convex deployment.
// The Go backend calls Convex query/mutation functions over HTTPS; Convex
// takes care of storage, pooling, and scaling, which replaces the direct
// PostgreSQL/Supabase connection entirely.
type ConvexCatalog struct {
	baseURL   string
	adminKey  string
	backendKey string
	client    *http.Client
}

type convexResult struct {
	Status       string          `json:"status"`
	Value        json.RawMessage `json:"value"`
	ErrorMessage string          `json:"errorMessage"`
}

// NewConvex creates a catalog bound to a Convex deployment. baseURL is the
// deployment URL (https://<deployment>.convex.cloud), adminKey is the Convex
// admin key used for privileged calls, and backendKey is the shared secret
// the Convex functions check before allowing writes.
func NewConvex(baseURL, adminKey, backendKey string) *ConvexCatalog {
	// Render hosts without IPv6 egress cannot reach some HTTPS endpoints;
	// reuse the shared IPv4-only transport from netutil when available.
	client := &http.Client{Timeout: 20 * time.Second}
	client.Transport = &http.Transport{DialContext: netutil.Ipv4Dialer().DialContext}
	return &ConvexCatalog{
		baseURL:    baseURL,
		adminKey:   adminKey,
		backendKey: backendKey,
		client:     client,
	}
}

// Close is a no-op because ConvexCatalog keeps no persistent connection.
func (c *ConvexCatalog) Close() error { return nil }

// Ping verifies connectivity and function availability by listing videos.
func (c *ConvexCatalog) Ping(ctx context.Context) error {
	_, err := c.listVideos(ctx, 1)
	return err
}

// Find returns the video identified by id, matching the legacy Supabase
// column layout so callers keep working with the same Item shape.
func (c *ConvexCatalog) Find(ctx context.Context, id string) (Item, error) {
	doc, err := c.callQuery(ctx, "videos:getVideo", map[string]any{"id": id})
	if err != nil {
		return Item{}, err
	}
	if doc == nil {
		return Item{}, ErrNotFound
	}
	return docToItem(doc)
}

// CanAccess reports whether the user may watch the video, honoring the
// free flag and the enrollments/subscriptions rules.
func (c *ConvexCatalog) CanAccess(ctx context.Context, userID string, item Item) (bool, error) {
	if item.Free {
		return true, nil
	}
	if item.CourseID == "" {
		return false, nil
	}
	out, err := c.callQuery(ctx, "videos:checkEnrollment", map[string]any{
		"studentId": userID,
		"courseId":  item.CourseID,
		"nowMs":     time.Now().UnixMilli(),
	})
	if err != nil {
		return false, err
	}
	var allowed bool
	if err := json.Unmarshal(out, &allowed); err != nil {
		return false, err
	}
	return allowed, nil
}

func (c *ConvexCatalog) listVideos(ctx context.Context, limit int) ([]json.RawMessage, error) {
	out, err := c.callQuery(ctx, "videos:listVideos", map[string]any{"limit": limit})
	if err != nil {
		return nil, err
	}
	var docs []json.RawMessage
	if err := json.Unmarshal(out, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (c *ConvexCatalog) callQuery(ctx context.Context, path string, args map[string]any) (json.RawMessage, error) {
	return c.call(ctx, "/api/query", path, args)
}

func (c *ConvexCatalog) callMutation(ctx context.Context, path string, args map[string]any) (json.RawMessage, error) {
	return c.call(ctx, "/api/mutation", path, args)
}

// CallQuery is an exported wrapper used by the Telegram bot package.
func (c *ConvexCatalog) CallQuery(ctx context.Context, path string, args map[string]any) (json.RawMessage, error) {
	return c.callQuery(ctx, path, args)
}

// CallMutation is an exported wrapper used by the Telegram bot package.
func (c *ConvexCatalog) CallMutation(ctx context.Context, path string, args map[string]any) (json.RawMessage, error) {
	return c.callMutation(ctx, path, args)
}

func (c *ConvexCatalog) call(ctx context.Context, endpoint, path string, args map[string]any) (json.RawMessage, error) {
	body, err := json.Marshal(map[string]any{"path": path, "args": args, "format": "json"})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.adminKey != "" {
		req.Header.Set("Authorization", "Convex "+c.adminKey)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	var result convexResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("convex: unexpected response: %w", err)
	}
	if result.Status != "success" {
		return nil, fmt.Errorf("convex %s: %s", path, result.ErrorMessage)
	}
	return result.Value, nil
}

// docToItem converts a Convex document into the catalog Item shape the rest
// of the backend expects. Missing fields default to values compatible with
// the legacy migration fallback.
func docToItem(doc json.RawMessage) (Item, error) {
	var raw map[string]any
	if err := json.Unmarshal(doc, &raw); err != nil {
		return Item{}, err
	}
	get := func(name string) string {
		v, ok := raw[name]
		if !ok || v == nil {
			return ""
		}
		switch s := v.(type) {
		case string:
			return s
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	getBool := func(name string) bool {
		v, ok := raw[name]
		if !ok || v == nil {
			return false
		}
		b, ok := v.(bool)
		return ok && b
	}
	status := get("status")
	published := getBool("is_published")
	if status == "" {
		status = "published"
		if published {
			published = true
		}
	}
	return Item{
		ID:         get("_id"),
		Title:      get("title"),
		CourseID:   get("course_id"),
		SourcePath: get("source_path"),
		YouTubeURL: get("youtube_url"),
		YouTubeID:  get("youtube_video_id"),
		MimeType:   orDefault(get("mime_type"), "video/mp4"),
		Status:     status,
		Published:  published,
		Free:       getBool("is_free"),
	}, nil
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

var ErrConvexUnavailable = errors.New("convex catalog unavailable")
