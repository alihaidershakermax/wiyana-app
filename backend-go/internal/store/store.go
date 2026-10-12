package store

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/netutil"
	"github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("not found")

type ViolationEvent struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	At       int64          `json:"at"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type UserState struct {
	Points     int              `json:"points"`
	Violations map[string]int   `json:"violations"`
	Events     []ViolationEvent `json:"events"`
	Ban        *Ban             `json:"ban,omitempty"`
}

type Ban struct {
	Reason       string `json:"reason"`
	Points       int    `json:"points"`
	BannedAt     int64  `json:"banned_at"`
	Until        int64  `json:"until_ts"`
	DurationSecs int    `json:"duration_seconds"`
}

func (b *Ban) Active(now time.Time) bool {
	return b != nil && b.Until > now.Unix()
}

type Grant struct {
	Token         string `json:"token"`
	UserID        string `json:"user_id"`
	VideoID       string `json:"video_id"`
	Mode          string `json:"mode"`
	Source        string `json:"source,omitempty"`
	CacheDir      string `json:"cache_dir,omitempty"`
	WatermarkText string `json:"watermark_text,omitempty"`
	DeviceID           string `json:"device_id,omitempty"`
	GoogleDriveSA      string `json:"-"`
	GoogleDriveAPIKey  string `json:"-"`
	CreatedAt          int64  `json:"created_at"`
	ExpiresAt     int64  `json:"expires_at"`
}

type Session struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	VideoID      string `json:"video_id"`
	CourseID     string `json:"course_id"`
	Platform     string `json:"platform,omitempty"`
	DisplayName  string `json:"display_name,omitempty"`
	TelegramID   string `json:"telegram_id,omitempty"`
	CreatedAt    int64  `json:"created_at"`
	ExpiresAt    int64  `json:"expires_at"`
	LastSeenAt   int64  `json:"last_seen_at"`
	RiskScore    int    `json:"risk_score"`
	Revoked      bool   `json:"revoked"`
	RevokeReason string `json:"revoke_reason,omitempty"`
	Token        string `json:"token,omitempty"`
}

type SecurityEvent struct {
	ID         string         `json:"id"`
	EventType  string         `json:"eventType"`
	UserID     string         `json:"userId"`
	VideoID    string         `json:"videoId,omitempty"`
	SessionID  string         `json:"sessionId,omitempty"`
	CourseID   string         `json:"courseId,omitempty"`
	Platform   string         `json:"platform,omitempty"`
	Confidence float64        `json:"confidence"`
	RiskScore  int            `json:"riskScore"`
	RiskLevel  string         `json:"riskLevel"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  int64          `json:"createdAt"`
}

type Store interface {
	Close() error
	LoadUserState(context.Context, string) (UserState, error)
	SaveUserState(context.Context, string, UserState, time.Duration) error
	PutGrant(context.Context, Grant, time.Duration) error
	GetGrant(context.Context, string) (Grant, error)
	DeleteGrant(context.Context, string) error
	PutSession(context.Context, Session, time.Duration) error
	GetSession(context.Context, string) (Session, error)
	DeleteSession(context.Context, string) error
	ListSessions(context.Context, string) ([]Session, error)
	PutEvent(context.Context, SecurityEvent, time.Duration) error
	ListEvents(context.Context, string, int) ([]SecurityEvent, error)
}

type RedisStore struct{ client *redis.Client }

func NewRedis(redisURL string) (*RedisStore, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	// IPv4-only connection: Render (Oregon) cannot reach external IPv6
	// addresses. Resolve the Redis host to an IPv4 literal first so the
	// driver never even sees an AAAA record.
	if h := options.Addr; h != "" {
		if newAddr, err := netutil.HostPortIPv4(context.Background(), h); err == nil {
			options.Addr = newAddr
		}
	}
	options.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return netutil.Ipv4Dialer().DialContext(ctx, network, addr)
	}
	client := redis.NewClient(options)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return &RedisStore{client: client}, nil
}

func (s *RedisStore) Close() error { return s.client.Close() }

func stateKey(id string) string   { return "wiyana:protection:state:" + id }
func grantKey(id string) string   { return "wiyana:protection:grant:" + id }
func sessionKey(id string) string { return "wiyana:playback:session:" + id }
func sessionSet(id string) string { return "wiyana:playback:user-sessions:" + id }
func eventKey(id string) string   { return "wiyana:security:event:" + id }
func eventList(id string) string  { return "wiyana:security:user-events:" + id }

func (s *RedisStore) LoadUserState(ctx context.Context, id string) (UserState, error) {
	data, err := s.client.Get(ctx, stateKey(id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return UserState{Violations: map[string]int{}}, nil
	}
	if err != nil {
		return UserState{}, err
	}
	var state UserState
	if err := json.Unmarshal(data, &state); err != nil {
		return UserState{}, err
	}
	if state.Violations == nil {
		state.Violations = map[string]int{}
	}
	return state, nil
}

func (s *RedisStore) SaveUserState(ctx context.Context, id string, state UserState, ttl time.Duration) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, stateKey(id), data, ttl).Err()
}

func (s *RedisStore) PutGrant(ctx context.Context, grant Grant, ttl time.Duration) error {
	data, err := json.Marshal(grant)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, grantKey(grant.Token), data, ttl).Err()
}

func (s *RedisStore) GetGrant(ctx context.Context, token string) (Grant, error) {
	data, err := s.client.Get(ctx, grantKey(token)).Bytes()
	if errors.Is(err, redis.Nil) {
		return Grant{}, ErrNotFound
	}
	if err != nil {
		return Grant{}, err
	}
	var grant Grant
	if err := json.Unmarshal(data, &grant); err != nil {
		return Grant{}, err
	}
	return grant, nil
}

func (s *RedisStore) DeleteGrant(ctx context.Context, token string) error {
	return s.client.Del(ctx, grantKey(token)).Err()
}

func (s *RedisStore) PutSession(ctx context.Context, session Session, ttl time.Duration) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	pipe := s.client.TxPipeline()
	pipe.Set(ctx, sessionKey(session.ID), data, ttl)
	pipe.SAdd(ctx, sessionSet(session.UserID), session.ID)
	pipe.Expire(ctx, sessionSet(session.UserID), ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (s *RedisStore) GetSession(ctx context.Context, id string) (Session, error) {
	data, err := s.client.Get(ctx, sessionKey(id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return Session{}, ErrNotFound
	}
	if err != nil {
		return Session{}, err
	}
	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return Session{}, err
	}
	return session, nil
}

func (s *RedisStore) DeleteSession(ctx context.Context, id string) error {
	return s.client.Del(ctx, sessionKey(id)).Err()
}

func (s *RedisStore) ListSessions(ctx context.Context, userID string) ([]Session, error) {
	ids, err := s.client.SMembers(ctx, sessionSet(userID)).Result()
	if err != nil {
		return nil, err
	}
	out := make([]Session, 0, len(ids))
	for _, id := range ids {
		session, err := s.GetSession(ctx, id)
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		out = append(out, session)
	}
	return out, nil
}

func (s *RedisStore) PutEvent(ctx context.Context, event SecurityEvent, ttl time.Duration) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	pipe := s.client.TxPipeline()
	pipe.Set(ctx, eventKey(event.ID), data, ttl)
	pipe.LPush(ctx, eventList(event.UserID), event.ID)
	pipe.LTrim(ctx, eventList(event.UserID), 0, 199)
	pipe.Expire(ctx, eventList(event.UserID), ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (s *RedisStore) ListEvents(ctx context.Context, userID string, limit int) ([]SecurityEvent, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	ids, err := s.client.LRange(ctx, eventList(userID), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	out := make([]SecurityEvent, 0, len(ids))
	for _, id := range ids {
		data, err := s.client.Get(ctx, eventKey(id)).Bytes()
		if errors.Is(err, redis.Nil) {
			continue
		}
		if err != nil {
			return nil, err
		}
		var event SecurityEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		out = append(out, event)
	}
	return out, nil
}

type MemoryStore struct {
	mu       sync.RWMutex
	states   map[string]UserState
	grants   map[string]Grant
	sessions map[string]Session
	events   map[string][]SecurityEvent
}

func NewMemory() *MemoryStore {
	return &MemoryStore{states: map[string]UserState{}, grants: map[string]Grant{}, sessions: map[string]Session{}, events: map[string][]SecurityEvent{}}
}
func (s *MemoryStore) Close() error { return nil }
func (s *MemoryStore) LoadUserState(_ context.Context, id string) (UserState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state := s.states[id]
	if state.Violations == nil {
		state.Violations = map[string]int{}
	}
	return state, nil
}
func (s *MemoryStore) SaveUserState(_ context.Context, id string, state UserState, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[id] = state
	return nil
}
func (s *MemoryStore) PutGrant(_ context.Context, grant Grant, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.grants[grant.Token] = grant
	return nil
}
func (s *MemoryStore) GetGrant(_ context.Context, token string) (Grant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	grant, ok := s.grants[token]
	if !ok || grant.ExpiresAt <= time.Now().Unix() {
		return Grant{}, ErrNotFound
	}
	return grant, nil
}
func (s *MemoryStore) DeleteGrant(_ context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.grants, token)
	return nil
}
func (s *MemoryStore) PutSession(_ context.Context, session Session, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID] = session
	return nil
}
func (s *MemoryStore) GetSession(_ context.Context, id string) (Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[id]
	if !ok || session.ExpiresAt <= time.Now().Unix() {
		return Session{}, ErrNotFound
	}
	return session, nil
}
func (s *MemoryStore) DeleteSession(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	return nil
}
func (s *MemoryStore) ListSessions(_ context.Context, userID string) ([]Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []Session{}
	for _, session := range s.sessions {
		if session.UserID == userID && session.ExpiresAt > time.Now().Unix() {
			out = append(out, session)
		}
	}
	return out, nil
}
func (s *MemoryStore) PutEvent(_ context.Context, event SecurityEvent, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := append([]SecurityEvent{event}, s.events[event.UserID]...)
	if len(items) > 200 {
		items = items[:200]
	}
	s.events[event.UserID] = items
	return nil
}
func (s *MemoryStore) ListEvents(_ context.Context, userID string, limit int) ([]SecurityEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := append([]SecurityEvent(nil), s.events[userID]...)
	if limit <= 0 || limit > len(items) {
		limit = len(items)
	}
	return items[:limit], nil
}
