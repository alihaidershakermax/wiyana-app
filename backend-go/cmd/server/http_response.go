package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/security"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/video"
)

func readJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "صيغة الطلب غير صالحة")
		return false
	}
	return true
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"success": false, "detail": message, "error": message})
}
func (s *server) securityError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, security.ErrBanned), errors.Is(err, security.ErrNoAccess):
		writeError(w, http.StatusForbidden, "لا تملك صلاحية تشغيل هذا الفيديو")
	case errors.Is(err, security.ErrSession), errors.Is(err, security.ErrGrant):
		writeError(w, http.StatusUnauthorized, "جلسة التشغيل أو المنحة غير صالحة")
	case errors.Is(err, video.ErrNotFound):
		writeError(w, http.StatusNotFound, "الفيديو غير موجود")
	default:
		writeError(w, http.StatusInternalServerError, "تعذر تنفيذ طلب التشغيل")
	}
}
func riskLevel(score int) string {
	if score >= 80 {
		return "CRITICAL"
	}
	if score >= 60 {
		return "HIGH_RISK"
	}
	if score >= 30 {
		return "SUSPICIOUS"
	}
	return "NORMAL"
}
