package security

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/config"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/store"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/video"
)

var (
	ErrBanned       = errors.New("user is banned")
	ErrNoAccess     = errors.New("user has no access to this video")
	ErrSession      = errors.New("playback session is invalid")
	ErrGrant        = errors.New("playback grant is invalid")
	ErrInvalidEvent = errors.New("unsupported security event")
)

type AccessStatus struct {
	Allowed           bool   `json:"allowed"`
	Banned            bool   `json:"banned"`
	ProtectionEnabled bool   `json:"protection_enabled"`
	Points            int    `json:"points"`
	Warning           bool   `json:"warning"`
	WarningThreshold  int    `json:"warning_threshold"`
	BanThreshold      int    `json:"ban_threshold"`
	BanReason         string `json:"ban_reason"`
	BanUntil          int64  `json:"ban_until"`
	BanDays           int    `json:"ban_days"`
}

type ReportResult struct {
	Success           bool   `json:"success"`
	ProtectionEnabled bool   `json:"protection_enabled"`
	UserID            string `json:"user_id"`
	ViolationType     string `json:"violation_type"`
	Points            int    `json:"points"`
	Status            string `json:"status"`
	Banned            bool   `json:"banned"`
	Behavioral        bool   `json:"behavioral"`
	WarningThreshold  int    `json:"warning_threshold"`
	BanThreshold      int    `json:"ban_threshold"`
	BanReason         string `json:"ban_reason"`
	BanUntil          int64  `json:"ban_until"`
	BanDays           int    `json:"ban_days"`
}

type PlaybackResponse struct {
	Success bool `json:"success"`
	Session struct {
		SessionID string `json:"sessionId"`
		UserID    string `json:"userId"`
		VideoID   string `json:"videoId"`
		CourseID  string `json:"courseId"`
		Platform  string `json:"platform,omitempty"`
		Status    string `json:"status"`
		ExpiresAt int64  `json:"expiresAt"`
	} `json:"session"`
	Token     string `json:"token"`
	Watermark struct {
		DisplayName string `json:"displayName"`
		TelegramID  string `json:"telegramId"`
		PrivacyID   string `json:"privacyId"`
		SessionID   string `json:"sessionId"`
	} `json:"watermark"`
}

type Service struct {
	cfg     config.Config
	store   store.Store
	catalog video.Catalog
}

func New(cfg config.Config, data store.Store, catalog video.Catalog) *Service {
	return &Service{cfg: cfg, store: data, catalog: catalog}
}

func (s *Service) Access(ctx context.Context, userID string) (AccessStatus, error) {
	state, err := s.store.LoadUserState(ctx, userID)
	if err != nil {
		return AccessStatus{}, err
	}
	status := AccessStatus{Allowed: s.cfg.ProtectionEnabled, ProtectionEnabled: s.cfg.ProtectionEnabled, Points: state.Points, Warning: state.Points >= s.cfg.WarningThreshold, WarningThreshold: s.cfg.WarningThreshold, BanThreshold: s.cfg.BanThreshold}
	if state.Ban != nil && state.Ban.Active(time.Now()) {
		status.Allowed = false
		status.Banned = true
		status.BanReason = state.Ban.Reason
		status.BanUntil = state.Ban.Until
		status.BanDays = max(1, int((state.Ban.Until-time.Now().Unix()+86399)/86400))
	}
	return status, nil
}

func (s *Service) ReportViolation(ctx context.Context, userID, violationType string, metadata map[string]any) (ReportResult, error) {
	if !s.cfg.ProtectionEnabled {
		return ReportResult{Success: true, ProtectionEnabled: false, UserID: userID, Status: "ok"}, nil
	}
	points := map[string]int{"screenshot_attempt": 10, "screen_recording": 10, "devtools_open": 5, "fullscreen_exit": 3}
	value, ok := points[violationType]
	if !ok {
		return ReportResult{}, ErrInvalidEvent
	}
	state, err := s.store.LoadUserState(ctx, userID)
	if err != nil {
		return ReportResult{}, err
	}
	if state.Violations == nil {
		state.Violations = map[string]int{}
	}
	now := time.Now()
	state.Points += value
	state.Violations[violationType]++
	state.Events = append(state.Events, store.ViolationEvent{ID: randomID(), Type: violationType, At: now.Unix(), Metadata: metadata})
	behavioral := 0
	cutoff := now.Unix() - int64(s.cfg.BehaviorWindowSeconds)
	trimmed := state.Events[:0]
	for _, event := range state.Events {
		if event.At >= cutoff {
			trimmed = append(trimmed, event)
			behavioral++
		}
	}
	state.Events = trimmed
	if behavioral >= s.cfg.EscalationViolations {
		state.Points += value
	}
	shouldBan := state.Points >= s.cfg.BanThreshold || (behavioral >= s.cfg.EscalationViolations)
	if shouldBan {
		state.Ban = &store.Ban{Reason: violationType, Points: state.Points, BannedAt: now.Unix(), Until: now.Add(time.Duration(s.cfg.BanDurationSeconds) * time.Second).Unix(), DurationSecs: s.cfg.BanDurationSeconds}
	}
	if err := s.store.SaveUserState(ctx, userID, state, time.Duration(max(s.cfg.BanDurationSeconds, s.cfg.BehaviorWindowSeconds)*2)*time.Second); err != nil {
		return ReportResult{}, err
	}
	result := ReportResult{Success: true, ProtectionEnabled: true, UserID: userID, ViolationType: violationType, Points: state.Points, Status: "ok", Banned: shouldBan, Behavioral: behavioral >= s.cfg.EscalationViolations, WarningThreshold: s.cfg.WarningThreshold, BanThreshold: s.cfg.BanThreshold}
	if shouldBan {
		result.Status = "banned"
		result.BanReason = state.Ban.Reason
		result.BanUntil = state.Ban.Until
		result.BanDays = max(1, int((state.Ban.Until-now.Unix()+86399)/86400))
	} else if state.Points >= s.cfg.WarningThreshold {
		result.Status = "warning"
	}
	return result, nil
}

func (s *Service) CreatePlaybackSession(ctx context.Context, userID, name, telegramID, videoID, courseID, platform string) (PlaybackResponse, error) {
	status, err := s.Access(ctx, userID)
	if err != nil {
		return PlaybackResponse{}, err
	}
	if status.Banned {
		return PlaybackResponse{}, ErrBanned
	}
	if !status.Allowed {
		return PlaybackResponse{}, ErrNoAccess
	}
	item, err := s.catalog.Find(ctx, videoID)
	if errors.Is(err, video.ErrNotFound) {
		return PlaybackResponse{}, video.ErrNotFound
	}
	if err != nil {
		return PlaybackResponse{}, err
	}
	if courseID == "" {
		courseID = item.CourseID
	}
	allowed, err := s.catalog.CanAccess(ctx, userID, item)
	if err != nil {
		return PlaybackResponse{}, err
	}
	if !allowed {
		return PlaybackResponse{}, ErrNoAccess
	}
	active, err := s.store.ListSessions(ctx, userID)
	if err != nil {
		return PlaybackResponse{}, err
	}
	for _, old := range active {
		if !old.Revoked {
			old.Revoked = true
			old.RevokeReason = "superseded_by_new_session"
			_ = s.store.PutSession(ctx, old, time.Until(time.Unix(old.ExpiresAt, 0)))
		}
	}
	now := time.Now()
	session := store.Session{ID: randomID(), UserID: userID, VideoID: videoID, CourseID: courseID, Platform: platform, DisplayName: name, TelegramID: telegramID, CreatedAt: now.Unix(), ExpiresAt: now.Add(time.Duration(s.cfg.SessionTTLSeconds) * time.Second).Unix(), LastSeenAt: now.Unix()}
	session.Token = s.signPlaybackToken(session)
	if err := s.store.PutSession(ctx, session, time.Duration(s.cfg.SessionTTLSeconds)*time.Second); err != nil {
		return PlaybackResponse{}, err
	}
	response := PlaybackResponse{Success: true, Token: session.Token}
	response.Session.SessionID, response.Session.UserID, response.Session.VideoID, response.Session.CourseID = session.ID, session.UserID, session.VideoID, session.CourseID
	response.Session.Platform, response.Session.Status, response.Session.ExpiresAt = session.Platform, "active", session.ExpiresAt
	response.Watermark.DisplayName, response.Watermark.TelegramID, response.Watermark.PrivacyID, response.Watermark.SessionID = name, telegramID, userID, session.ID[:min(12, len(session.ID))]
	return response, nil
}

func (s *Service) Heartbeat(ctx context.Context, userID, sessionID, videoID string, risk int) (store.Session, error) {
	session, err := s.store.GetSession(ctx, sessionID)
	if err != nil {
		return store.Session{}, ErrSession
	}
	if session.UserID != userID || session.VideoID != videoID || session.Revoked || session.ExpiresAt <= time.Now().Unix() {
		return store.Session{}, ErrSession
	}
	if risk >= 80 {
		session.Revoked = true
		session.RevokeReason = "risk_score_threshold"
		_ = s.store.PutSession(ctx, session, time.Until(time.Unix(session.ExpiresAt, 0)))
		return store.Session{}, ErrSession
	}
	session.LastSeenAt, session.RiskScore = time.Now().Unix(), risk
	if err := s.store.PutSession(ctx, session, time.Until(time.Unix(session.ExpiresAt, 0))); err != nil {
		return store.Session{}, err
	}
	return session, nil
}

func (s *Service) ValidatePlaybackToken(ctx context.Context, token, userID, videoID string) (store.Session, error) {
	claims, err := parsePlaybackToken(token, s.cfg.JWTSecret)
	if err != nil {
		return store.Session{}, ErrGrant
	}
	if claims.UserID != userID || claims.VideoID != videoID || claims.ExpiresAt <= time.Now().Unix() {
		return store.Session{}, ErrGrant
	}
	session, err := s.store.GetSession(ctx, claims.SessionID)
	if err != nil || session.Token != token || session.Revoked || session.ExpiresAt <= time.Now().Unix() {
		return store.Session{}, ErrGrant
	}
	return session, nil
}

func (s *Service) CreateGrant(ctx context.Context, userID, videoID, playbackToken string) (store.Grant, error) {
	return s.CreateGrantForDevice(ctx, userID, videoID, playbackToken, "")
}

// CreateGrantForDevice binds the playback grant to a device fingerprint so the
// same link cannot be replayed from another machine or downloaded to disk and
// played with a generic player (the key file is served only to the bound
// device). When deviceID is empty the binding is skipped (e.g. desktop web).
func (s *Service) CreateGrantForDevice(ctx context.Context, userID, videoID, playbackToken, deviceID string) (store.Grant, error) {
	item, err := s.catalog.Find(ctx, videoID)
	if err != nil {
		return store.Grant{}, err
	}
	if playbackToken == "" {
		return store.Grant{}, ErrGrant
	}
	session, err := s.ValidatePlaybackToken(ctx, playbackToken, userID, videoID)
	if err != nil {
		return store.Grant{}, err
	}
	allowed, err := s.catalog.CanAccess(ctx, userID, item)
	if err != nil || !allowed {
		return store.Grant{}, ErrNoAccess
	}
	mode, source := "hls", item.SourcePath
	if source == "" {
		if item.YouTubeURL != "" {
			mode, source = "direct", item.YouTubeURL
		} else {
			return store.Grant{}, errors.New("video has no playable source")
		}
	}
	if mode == "hls" {
		safeSource, err := video.SafePath(s.cfg.UploadDir, source)
		if err != nil {
			return store.Grant{}, fmt.Errorf("invalid video source: %w", err)
		}
		source = safeSource
	}
	now := time.Now()
	grant := store.Grant{Token: randomID(), UserID: userID, VideoID: videoID, Mode: mode, Source: source, WatermarkText: fmt.Sprintf("%s | %s | %s", session.DisplayName, session.TelegramID, session.ID[:min(12, len(session.ID))]), DeviceID: deviceID, CreatedAt: now.Unix(), ExpiresAt: now.Add(time.Duration(s.cfg.GrantTTLSeconds) * time.Second).Unix()}
	if strings.HasPrefix(source, "gdrive://") {
		grant.GoogleDriveSA = s.cfg.GoogleDriveSA
		grant.GoogleDriveAPIKey = s.cfg.GoogleDriveAPIKey
	}
	grant.CacheDir = fmt.Sprintf("%s/%s/%s", strings.TrimRight(s.cfg.HLSDir, "/"), userID, videoID)
	if err := s.store.PutGrant(ctx, grant, time.Duration(s.cfg.GrantTTLSeconds)*time.Second); err != nil {
		return store.Grant{}, err
	}
	_ = session
	return grant, nil
}

func (s *Service) GetGrant(ctx context.Context, grantToken string) (store.Grant, error) {
	grant, err := s.store.GetGrant(ctx, grantToken)
	if err != nil || grant.ExpiresAt <= time.Now().Unix() {
		return store.Grant{}, ErrGrant
	}
	return grant, nil
}

// GetGrantForDevice returns the grant only if the caller's device fingerprint
// matches the one the grant was issued for. A link captured by screen recording
// or file sharing keeps working for the key file only on the original device.
func (s *Service) GetGrantForDevice(ctx context.Context, grantToken, deviceID string) (store.Grant, error) {
	grant, err := s.GetGrant(ctx, grantToken)
	if err != nil {
		return store.Grant{}, err
	}
	if grant.DeviceID != "" && grant.DeviceID != deviceID {
		return store.Grant{}, ErrGrant
	}
	return grant, nil
}
func (s *Service) RevokeSession(ctx context.Context, sessionID, reason string) error {
	session, err := s.store.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}
	session.Revoked = true
	session.RevokeReason = reason
	return s.store.PutSession(ctx, session, time.Until(time.Unix(session.ExpiresAt, 0)))
}
func (s *Service) ListSessions(ctx context.Context, userID string) ([]store.Session, error) {
	return s.store.ListSessions(ctx, userID)
}
func (s *Service) Events(ctx context.Context, userID string, limit int) ([]store.SecurityEvent, error) {
	return s.store.ListEvents(ctx, userID, limit)
}

func (s *Service) Telemetry(ctx context.Context, userID, eventType, videoID, sessionID, courseID, platform string, confidence float64, metadata map[string]any) (store.SecurityEvent, error) {
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}
	risk := 0
	switch eventType {
	case "screenshot_attempt", "screen_recording":
		risk = 80
	case "devtools_open":
		risk = 60
	case "fullscreen_exit":
		risk = 20
	case "token_anomaly", "invalid_playback_token":
		risk = 90
	default:
		risk = 10
	}
	risk = int(float64(risk) * confidence)
	level := "NORMAL"
	if risk >= 80 {
		level = "CRITICAL"
	} else if risk >= 60 {
		level = "HIGH_RISK"
	} else if risk >= 30 {
		level = "SUSPICIOUS"
	}
	event := store.SecurityEvent{ID: randomID(), EventType: eventType, UserID: userID, VideoID: videoID, SessionID: sessionID, CourseID: courseID, Platform: platform, Confidence: confidence, RiskScore: risk, RiskLevel: level, Metadata: metadata, CreatedAt: time.Now().Unix()}
	if err := s.store.PutEvent(ctx, event, 30*24*time.Hour); err != nil {
		return store.SecurityEvent{}, err
	}
	return event, nil
}

func (s *Service) signPlaybackToken(session store.Session) string {
	payload, _ := json.Marshal(map[string]any{"sid": session.ID, "uid": session.UserID, "vid": session.VideoID, "exp": session.ExpiresAt})
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, []byte(s.cfg.JWTSecret))
	_, _ = mac.Write([]byte(encoded))
	return encoded + "." + hex.EncodeToString(mac.Sum(nil))
}

type playbackClaims struct {
	SessionID, UserID, VideoID string
	ExpiresAt                  int64
}

func parsePlaybackToken(token, secret string) (playbackClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return playbackClaims{}, ErrGrant
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(parts[0]))
	if !hmac.Equal([]byte(parts[1]), []byte(hex.EncodeToString(mac.Sum(nil)))) {
		return playbackClaims{}, ErrGrant
	}
	data, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return playbackClaims{}, ErrGrant
	}
	var raw struct {
		SID string `json:"sid"`
		UID string `json:"uid"`
		VID string `json:"vid"`
		Exp int64  `json:"exp"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return playbackClaims{}, ErrGrant
	}
	return playbackClaims{raw.SID, raw.UID, raw.VID, raw.Exp}, nil
}
func randomID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
