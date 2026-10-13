package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/auth"
)

func (s *server) loginTelegram(w http.ResponseWriter, r *http.Request) {
	var body struct {
		InitData string `json:"init_data"`
	}
	if !readJSON(w, r, &body) {
		return
	}
	identity, err := auth.VerifyTelegramInitData(body.InitData, s.cfg.BotToken, s.cfg.Environment)
	if err != nil {
		// TEMP DEBUG: detailed reason for diagnosis
		reason := "unknown"
		if body.InitData == "" {
			reason = "empty_init_data"
		} else if len(body.InitData) < 40 {
			reason = "init_data_too_short"
		} else {
			reason = "signature_or_expired"
		}
		writeError(w, http.StatusUnauthorized, "بيانات Telegram غير صالحة [debug:"+reason+"]")
		return
	}
	// Owners listed in BOT_ADMIN_IDS sign in with the owner role so the
	// admin panel is reachable in production (Telegram login otherwise
	// yields the plain "student" role).
	for _, adminID := range s.cfg.BotAdminIDs {
		if strconv.FormatInt(adminID, 10) == identity.TelegramID {
			identity.Role = "owner"
			break
		}
	}
	// Persist the user in SQLite (keeps the existing role on re-login).
	if s.adb != nil {
		if tid, perr := strconv.ParseInt(identity.TelegramID, 10, 64); perr == nil {
			if u, uerr := s.adb.UpsertByTelegram(tid, "", identity.Name, "", "", identity.Role); uerr == nil && u != nil {
				identity.ID = u.ID
				if strings.ToLower(u.Role) == "owner" || strings.ToLower(u.Role) == "admin" {
					identity.Role = strings.ToLower(u.Role)
				}
			}
		}
	}
	token, err := auth.SignJWT(identity, s.cfg.JWTSecret, 30*24*time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "تعذر إصدار جلسة الدخول")
		return
	}
	refresh, err := auth.SignJWT(identity, s.cfg.JWTSecret, 60*24*time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "تعذر إصدار جلسة الدخول")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"access_token": token, "refresh_token": refresh, "token_type": "bearer", "user": identity})
}
func (s *server) identity(w http.ResponseWriter, r *http.Request) (auth.Identity, int) {
	identity, err := auth.Resolve(r, auth.Config{JWTSecret: s.cfg.JWTSecret, BotToken: s.cfg.BotToken, Environment: s.cfg.Environment})
	if err == nil {
		return identity, 0
	}
	if s.cfg.AllowDevAuth && !strings.EqualFold(s.cfg.Environment, "production") {
		return auth.Identity{ID: "dev-user", TelegramID: "dev-user", Name: "مستخدم الاختبار", Role: "owner"}, 0
	}
	writeError(w, http.StatusUnauthorized, "يلزم تسجيل الدخول")
	return auth.Identity{}, http.StatusUnauthorized
}
func (s *server) admin(w http.ResponseWriter, r *http.Request) (auth.Identity, bool) {
	identity, status := s.identity(w, r)
	if status != 0 {
		return auth.Identity{}, false
	}
	role := strings.ToLower(identity.Role)
	if role != "admin" && role != "owner" && role != "teacher" {
		writeError(w, http.StatusForbidden, "لا تملك صلاحية الإدارة")
		return auth.Identity{}, false
	}
	return identity, true
}
