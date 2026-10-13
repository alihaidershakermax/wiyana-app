package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// One-time auth tokens for fallback when Telegram initData is unavailable.
// Token -> {telegramID, expiry}. Single-use, 10min TTL.
var oneTimeTokens = struct {
	sync.Mutex
	tokens map[string]oneTimeToken
}{tokens: make(map[string]oneTimeToken)}

type oneTimeToken struct {
	TelegramID int64
	Expiry     time.Time
}

// CreateOneTimeToken generates a single-use auth token for a Telegram user.
// Valid for 10 minutes.
func CreateOneTimeToken(telegramID int64) string {
	oneTimeTokens.Lock()
	defer oneTimeTokens.Unlock()
	now := time.Now()
	for k, v := range oneTimeTokens.tokens {
		if now.After(v.Expiry) {
			delete(oneTimeTokens.tokens, k)
		}
	}
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	token := "wot_" + hex.EncodeToString(b)
	oneTimeTokens.tokens[token] = oneTimeToken{TelegramID: telegramID, Expiry: now.Add(10 * time.Minute)}
	return token
}

// ConsumeOneTimeToken validates and consumes a one-time token.
// Returns the Telegram user ID and true if valid, 0 and false otherwise.
func ConsumeOneTimeToken(token string) (int64, bool) {
	oneTimeTokens.Lock()
	defer oneTimeTokens.Unlock()
	v, ok := oneTimeTokens.tokens[token]
	if !ok || time.Now().After(v.Expiry) {
		return 0, false
	}
	delete(oneTimeTokens.tokens, token)
	return v.TelegramID, true
}
