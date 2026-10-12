package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type cachedHTTPResponse struct {
	Status int    `json:"status"`
	ETag   string `json:"etag"`
	Body   []byte `json:"body"`
}

type cacheResponseWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (s *server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			s.setCORS(w, r)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		clientIP := requestIP(r)
		if !s.limiter.allow(clientIP+":"+r.URL.Path, s.cfg.HTTPRateLimit, time.Minute) {
			writeError(w, http.StatusTooManyRequests, "تم تجاوز عدد الطلبات المسموح به")
			return
		}
		s.setCORS(w, r)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), display-capture=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'; object-src 'none'; base-uri 'none'")
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
		handler := http.Handler(next)

		if s.cfg.CacheEnabled && s.cache != nil {
			handler = s.cacheMiddleware(handler)
		}
		handler.ServeHTTP(w, r)

	})
}
func (w *cacheResponseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (w *cacheResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	if w.status >= 200 && w.status < 300 && w.body.Len() <= 2<<20 {
		_, _ = w.body.Write(data)
	}
	return w.ResponseWriter.Write(data)
}
func (s *server) cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isCacheableRequest(r) {
			next.ServeHTTP(w, r)
			if r.Method != http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v1/") {
				cacheCtx, cancel := context.WithTimeout(context.Background(), time.Second)
				_ = s.cache.DeletePrefix(cacheCtx, "http:")
				_ = s.cache.DeletePrefix(cacheCtx, "catalog:")
				cancel()
			}
			return
		}
		key := cacheKey(r)
		if raw, ok, err := s.cache.Get(r.Context(), key); err == nil && ok {
			var cached cachedHTTPResponse
			if json.Unmarshal(raw, &cached) == nil && cached.Status == http.StatusOK {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.Header().Set("Cache-Control", fmt.Sprintf("private, max-age=%d, stale-while-revalidate=%d", s.cfg.CacheTTLSeconds, s.cfg.CacheTTLSeconds))
				w.Header().Set("ETag", cached.ETag)
				w.Header().Set("X-Cache", "HIT")
				if r.Header.Get("If-None-Match") == cached.ETag {
					w.WriteHeader(http.StatusNotModified)
					return
				}
				w.WriteHeader(cached.Status)
				_, _ = w.Write(cached.Body)
				return
			}
		}
		w.Header().Set("X-Cache", "MISS")
		captured := &cacheResponseWriter{ResponseWriter: w}
		next.ServeHTTP(captured, r)
		if captured.status != http.StatusOK || captured.body.Len() == 0 || captured.body.Len() > 2<<20 {
			return
		}
		body := captured.body.Bytes()
		sum := sha256.Sum256(body)
		cached := cachedHTTPResponse{Status: captured.status, ETag: `"` + hex.EncodeToString(sum[:]) + `"`, Body: append([]byte(nil), body...)}
		encoded, err := json.Marshal(cached)
		if err == nil {
			cacheCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			_ = s.cache.Set(cacheCtx, key, encoded, time.Duration(s.cfg.CacheTTLSeconds)*time.Second)
			cancel()
		}
	})
}
func isCacheableRequest(r *http.Request) bool {
	if r.Method != http.MethodGet || !strings.HasPrefix(r.URL.Path, "/api/v1/") {
		return false
	}
	path := r.URL.Path
	return !strings.HasPrefix(path, "/api/v1/auth/") &&
		!strings.HasPrefix(path, "/api/v1/playback/") &&
		!strings.HasPrefix(path, "/api/v1/security/")
}
func cacheKey(r *http.Request) string {
	authDigest := sha256.Sum256([]byte(r.Header.Get("Authorization") + "|" + r.Header.Get("X-Telegram-Init-Data")))
	return "http:" + r.Method + ":" + r.URL.RequestURI() + ":user:" + hex.EncodeToString(authDigest[:])
}
func (s *server) setCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	for _, allowed := range s.cfg.AllowedOrigins {
		if subtle.ConstantTimeCompare([]byte(origin), []byte(allowed)) == 1 {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Telegram-Init-Data, X-Playback-Token, X-Device-ID")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			return
		}
	}
}
func requestIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
