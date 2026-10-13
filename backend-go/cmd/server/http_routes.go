package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/cache"
)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/health" || r.URL.Path == "/api/v1/health" {
		stats := cache.Stats{}
		if s.cache != nil {
			stats = s.cache.Stats()
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "wiyana-go", "time": time.Now().UTC(), "cache": stats})
		return
	}
	if r.URL.Path == "/api/v1/auth/telegram-login" && r.Method == http.MethodPost {
		s.loginTelegram(w, r)
		return
	}
	if r.URL.Path == "/api/v1/auth/token-login" && r.Method == http.MethodPost {
		s.loginViaToken(w, r)
		return
	}
	// TEMP DEBUG: view recent login attempt metadata
	if r.URL.Path == "/api/v1/auth/debug-log" && r.Method == http.MethodGet {
		loginDebugLog.Lock()
		entries := make([]map[string]any, len(loginDebugLog.entries))
		copy(entries, loginDebugLog.entries)
		loginDebugLog.Unlock()
		writeJSON(w, http.StatusOK, map[string]any{"entries": entries})
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/v1/auth/") && s.authCompatRoute(w, r) {
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/stream/hls/") {
		s.streamRoute(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/stream/direct/") && r.Method == http.MethodGet {
		s.playbackDirect(w, r)
		return
	}
	if r.URL.Path == "/api/check-access" && r.Method == http.MethodGet {
		s.access(w, r)
		return
	}
	if r.URL.Path == "/api/report-violation" && r.Method == http.MethodPost {
		s.reportViolation(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/v1/playback/") {
		s.playbackRoute(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/v1/security/") {
		s.securityRoute(w, r)
		return
	}
	if r.URL.Path == "/api/v1/auth/me" && r.Method == http.MethodGet {
		identity, status := s.identity(w, r)
		if status != 0 {
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "user": identity})
		return
	}
	if s.adminPanelRoute(w, r) {
		return
	}
	if s.compatRoute(w, r) {
		return
	}
	if s.spa != nil && (r.Method == http.MethodGet || r.Method == http.MethodHead) {
		s.spa.ServeHTTP(w, r)
		return
	}
	writeError(w, http.StatusNotFound, "المسار غير موجود")
}
