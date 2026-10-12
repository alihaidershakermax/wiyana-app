package main

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/auth"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/security"
)

func (s *server) access(w http.ResponseWriter, r *http.Request) {
	identity, status := s.identity(w, r)
	if status != 0 {
		return
	}
	result, err := s.security.Access(r.Context(), identity.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "تعذر قراءة حالة الحماية")
		return
	}
	writeJSON(w, http.StatusOK, result)
}
func (s *server) reportViolation(w http.ResponseWriter, r *http.Request) {
	identity, status := s.identity(w, r)
	if status != 0 {
		return
	}
	var body struct {
		ViolationType string         `json:"violation_type"`
		LessonID      string         `json:"lesson_id"`
		Meta          map[string]any `json:"meta"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	if body.Meta == nil {
		body.Meta = map[string]any{}
	}
	if body.LessonID != "" {
		body.Meta["lesson_id"] = body.LessonID
	}
	result, err := s.security.ReportViolation(r.Context(), identity.ID, body.ViolationType, body.Meta)
	if err != nil {
		if errors.Is(err, security.ErrInvalidEvent) {
			writeError(w, http.StatusBadRequest, "نوع حدث أمني غير معروف")
		} else {
			writeError(w, http.StatusInternalServerError, "تعذر تسجيل الحدث الأمني")
		}
		return
	}
	writeJSON(w, http.StatusOK, result)
}
func (s *server) playbackRoute(w http.ResponseWriter, r *http.Request) {
	identity, status := s.identity(w, r)
	if status != 0 {
		return
	}
	switch {
	case r.URL.Path == "/api/v1/playback/session" && r.Method == http.MethodPost:
		var body struct {
			VideoID  string `json:"videoId"`
			CourseID string `json:"courseId"`
			Platform string `json:"platform"`
		}
		if !readJSON(w, r, &body) {
			return
		}
		response, err := s.security.CreatePlaybackSession(r.Context(), identity.ID, identity.Name, identity.TelegramID, body.VideoID, body.CourseID, body.Platform)
		if err != nil {
			s.securityError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, response)
	case r.URL.Path == "/api/v1/playback/heartbeat" && r.Method == http.MethodPost:
		var body struct {
			SessionID string `json:"sessionId"`
			VideoID   string `json:"videoId"`
			RiskScore int    `json:"riskScore"`
		}
		if !readJSON(w, r, &body) {
			return
		}
		_, err := s.security.Heartbeat(r.Context(), identity.ID, body.SessionID, body.VideoID, body.RiskScore)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"success": false, "status": "revoked"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "status": "active", "riskScore": body.RiskScore, "riskLevel": riskLevel(body.RiskScore)})
	default:
		writeError(w, http.StatusNotFound, "مسار التشغيل غير موجود")
	}
}
func (s *server) streamRoute(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/stream/hls/")
	parts := strings.Split(rest, "/")
	if len(parts) == 0 || parts[0] == "" {
		writeError(w, http.StatusNotFound, "معرّف الفيديو مفقود")
		return
	}
	videoID, err := url.PathUnescape(parts[0])
	if err != nil {
		writeError(w, http.StatusBadRequest, "معرّف فيديو غير صالح")
		return
	}
	identity, status := s.identity(w, r)
	if status != 0 {
		return
	}
	if r.Method == http.MethodPost && len(parts) == 1 {
		token := r.Header.Get("X-Playback-Token")
		deviceID := r.Header.Get("X-Device-ID")
		grant, err := s.security.CreateGrantForDevice(r.Context(), identity.ID, videoID, token, deviceID)
		if err != nil {
			s.securityError(w, err)
			return
		}
		playlist := fmt.Sprintf("/api/stream/hls/%s/playlist.m3u8?grant=%s", url.PathEscape(videoID), url.QueryEscape(grant.Token))
		if grant.Mode == "direct" {
			playlist = fmt.Sprintf("/api/stream/direct/%s?grant=%s", url.PathEscape(videoID), url.QueryEscape(grant.Token))
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "mode": grant.Mode, "video_id": videoID, "playlist_url": playlist, "expires_in": grant.ExpiresAt - time.Now().Unix()})
		return
	}
	if r.Method == http.MethodGet && len(parts) >= 2 {
		s.serveHLS(w, r, identity.ID, videoID, parts[1])
		return
	}
	writeError(w, http.StatusMethodNotAllowed, "طريقة الطلب غير مدعومة")
}
func (s *server) serveHLS(w http.ResponseWriter, r *http.Request, userID, videoID, file string) {
	grantToken := r.URL.Query().Get("grant")
	deviceID := r.Header.Get("X-Device-ID")
	grant, err := s.security.GetGrantForDevice(r.Context(), grantToken, deviceID)
	if err != nil || grant.UserID != userID || grant.VideoID != videoID || grant.Mode != "hls" {
		writeError(w, http.StatusForbidden, "منحة التشغيل غير صالحة")
		return
	}
	base := filepath.Clean(grant.CacheDir)
	name, err := url.PathUnescape(file)
	if err != nil || name == "" || filepath.Base(name) != name || strings.Contains(name, "..") {
		writeError(w, http.StatusForbidden, "مسار ملف غير صالح")
		return
	}
	if _, err := s.streamer.Generate(r.Context(), grant); err != nil {
		writeError(w, http.StatusServiceUnavailable, "تعذر تجهيز البث")
		return
	}
	target := filepath.Join(base, name)
	rel, err := filepath.Rel(base, target)
	if err != nil || strings.HasPrefix(rel, "..") {
		writeError(w, http.StatusForbidden, "مسار ملف غير صالح")
		return
	}
	if name == "playlist.m3u8" {
		data, err := os.ReadFile(target)
		if err != nil {
			writeError(w, http.StatusNotFound, "قائمة التشغيل غير موجودة")
			return
		}
		// Rewrite every referenced file (segments + stream.key) to route
		// through this endpoint with the same grant token, so nothing is ever
		// served directly from the filesystem and every request is validated.
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			// Rewrite the #EXT-X-KEY URI to route through this endpoint too,
			// so a copied playlist can never decrypt content outside this
			// route's grant and device checks.
			if strings.HasPrefix(line, "#EXT-X-KEY:") {
				old := extractKeyURI(line)
				if old != "" {
					q := "?grant=" + url.QueryEscape(grant.Token)
					if grant.DeviceID != "" {
						q += "&device=" + url.QueryEscape(grant.DeviceID)
					}
					lines[i] = "#EXT-X-KEY:METHOD=AES-128,URI=\"/api/stream/hls/" + url.PathEscape(videoID) + "/stream.key" + q + "\",IV=" + extractKeyIV(line)
				}
				continue
			}
			if line != "" && !strings.HasPrefix(line, "#") {
				base := filepath.Base(line)
				q := "?grant=" + url.QueryEscape(grant.Token)
				if grant.DeviceID != "" {
					q += "&device=" + url.QueryEscape(grant.DeviceID)
				}
				lines[i] = fmt.Sprintf("/api/stream/hls/%s/%s%s", url.PathEscape(videoID), url.PathEscape(base), q)
			}
		}
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Header().Set("Cache-Control", "private, no-store, max-age=0, must-revalidate")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write([]byte(strings.Join(lines, "\n")))

		return
	}
	info, err := os.Stat(target)
	if err != nil || info.IsDir() {
		writeError(w, http.StatusNotFound, "مقطع البث غير موجود")
		return
	}
	// The key file must be requested on the route itself (never directly), so
	// a shared link without this route's grant check cannot decrypt content.
	if name == "stream.key" {
		key, err := os.ReadFile(target)
		if err != nil {
			writeError(w, http.StatusNotFound, "مفتاح التشفير غير متاح")
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Cache-Control", "private, no-store, max-age=0, must-revalidate")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Length", strconv.Itoa(len(key)))
		_, _ = w.Write(key)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(name)))
	w.Header().Set("Cache-Control", "private, no-store, max-age=0, must-revalidate")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeFile(w, r, target)
}
// extractKeyURI parses the URI value out of a #EXT-X-KEY line, keeping any
// legacy (build directory) path the client-side player would otherwise fail on.
func extractKeyURI(line string) string {
	start := strings.Index(line, `URI="`)
	if start < 0 {
		return ""
	}
	start += 5
	end := strings.Index(line[start:], `"`)
	if end < 0 {
		return line[start:]
	}
	return line[start : start+end]
}
func extractKeyIV(line string) string {
	idx := strings.Index(line, ",IV=")
	if idx < 0 {
		return "00000000000000000000000000000000"
	}
	return line[idx+4:]
}
func (s *server) playbackDirect(w http.ResponseWriter, r *http.Request) {
	identity, status := s.identity(w, r)
	if status != 0 {
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, "/api/stream/direct/")
	videoID, _ := url.PathUnescape(strings.Trim(rest, "/"))
	grant, err := s.security.GetGrant(r.Context(), r.URL.Query().Get("grant"))
	if err != nil || grant.UserID != identity.ID || grant.VideoID != videoID || grant.Mode != "direct" {
		writeError(w, http.StatusForbidden, "منحة التشغيل غير صالحة")
		return
	}
	http.Redirect(w, r, grant.Source, http.StatusTemporaryRedirect)
}
func (s *server) securityRoute(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/v1/security/admin/") {
		identity, ok := s.admin(w, r)
		if !ok {
			return
		}
		s.adminRoute(w, r, identity)
		return
	}
	identity, status := s.identity(w, r)
	if status != 0 {
		return
	}
	if r.URL.Path == "/api/v1/security/telemetry" && r.Method == http.MethodPost {
		var body struct {
			EventType  string         `json:"eventType"`
			VideoID    string         `json:"videoId"`
			SessionID  string         `json:"sessionId"`
			Confidence float64        `json:"confidence"`
			Metadata   map[string]any `json:"metadata"`
			CourseID   string         `json:"courseId"`
			Platform   string         `json:"platform"`
		}
		if !readJSON(w, r, &body) {
			return
		}
		event, err := s.security.Telemetry(r.Context(), identity.ID, body.EventType, body.VideoID, body.SessionID, body.CourseID, body.Platform, body.Confidence, body.Metadata)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "تعذر تسجيل القياس الأمني")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "eventId": event.ID, "riskScore": event.RiskScore, "riskLevel": event.RiskLevel})
		return
	}
	writeError(w, http.StatusNotFound, "مسار الأمان غير موجود")
}
func (s *server) adminRoute(w http.ResponseWriter, r *http.Request, identity auth.Identity) {
	switch {
	case r.URL.Path == "/api/v1/security/admin/events" && r.Method == http.MethodGet:
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		events, err := s.security.Events(r.Context(), identity.ID, limit)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "تعذر قراءة الأحداث")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "events": events})
	case r.URL.Path == "/api/v1/security/admin/sessions" && r.Method == http.MethodGet:
		sessions, err := s.security.ListSessions(r.Context(), identity.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "تعذر قراءة الجلسات")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "sessions": sessions})
	case r.URL.Path == "/api/v1/security/admin/revoke-session" && r.Method == http.MethodPost:
		var body struct {
			SessionID string `json:"sessionId"`
			Reason    string `json:"reason"`
		}
		if !readJSON(w, r, &body) {
			return
		}
		if body.Reason == "" {
			body.Reason = "admin_action"
		}
		if err := s.security.RevokeSession(r.Context(), body.SessionID, body.Reason); err != nil {
			writeError(w, http.StatusNotFound, "الجلسة غير موجودة")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true})
	case r.URL.Path == "/api/v1/security/admin/revoke-all" && r.Method == http.MethodPost:
		var body struct {
			UserID string `json:"userId"`
		}
		if !readJSON(w, r, &body) {
			return
		}
		sessions, _ := s.security.ListSessions(r.Context(), body.UserID)
		for _, session := range sessions {
			_ = s.security.RevokeSession(r.Context(), session.ID, "admin_revoke_all")
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true})
	default:
		if strings.HasPrefix(r.URL.Path, "/api/v1/security/admin/users/") && strings.HasSuffix(r.URL.Path, "/history") && r.Method == http.MethodGet {
			raw := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/security/admin/users/"), "/history")
			userID, _ := url.PathUnescape(raw)
			limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
			events, err := s.security.Events(r.Context(), userID, limit)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "تعذر قراءة التاريخ الأمني")
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"success": true, "events": events})
			return
		}
		writeError(w, http.StatusNotFound, "مسار الإدارة غير موجود")
	}
}

