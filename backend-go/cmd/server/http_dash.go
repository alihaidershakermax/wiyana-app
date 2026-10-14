package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/auth"
)

// ClearKey license server for DASH encrypted videos.
// dash.js sends: {"kids": ["base64url_key_id", ...]}
// We return: {"keys": [{"kty":"oct","kid":"...","k":"..."}]}
// Keys are base64url-encoded (no padding) per ClearKey spec.

// In-memory key store: key_id (hex) -> key (hex).
// For Phase 1, loaded from DASH_KEYS env: "keyid1:key1,keyid2:key2"
// TODO: move to DB for per-video keys.
var (
	dashKeys   = make(map[string]string)
	dashKeysMu sync.RWMutex
)

func initDashKeys() {
	// Load keys from DASH_KEYS env: "keyid1:key1,keyid2:key2" (hex)
	if envKeys := os.Getenv("DASH_KEYS"); envKeys != "" {
		for _, pair := range strings.Split(envKeys, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) == 2 {
				setDashKey(parts[0], parts[1])
			}
		}
	}
}

func setDashKey(keyIDHex, keyHex string) {
	dashKeysMu.Lock()
	defer dashKeysMu.Unlock()
	dashKeys[strings.ToLower(keyIDHex)] = keyHex
}

func getDashKey(keyIDHex string) (string, bool) {
	dashKeysMu.RLock()
	defer dashKeysMu.RUnlock()
	k, ok := dashKeys[strings.ToLower(keyIDHex)]
	return k, ok
}

// hexToBase64URL converts hex string to base64url (no padding)
func hexToBase64URL(h string) (string, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// base64URLToHex converts base64url (no padding) to hex
func base64URLToHex(s string) (string, error) {
	// Add padding if needed
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type clearKeyRequest struct {
	Kids []string `json:"kids"`
	Type string   `json:"type,omitempty"`
}

type clearKeyKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	K   string `json:"k"`
}

type clearKeyResponse struct {
	Keys []clearKeyKey `json:"keys"`
	Type string        `json:"type,omitempty"`
}

// handleDashLicense is the ClearKey license server endpoint.
// Requires valid JWT (user must be logged in). Optionally checks enrollment.
func (s *server) handleDashLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify JWT — user must be authenticated
	identity, err := auth.Resolve(r, auth.Config{JWTSecret: s.cfg.JWTSecret, BotToken: s.cfg.BotToken, Environment: s.cfg.Environment})
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	_ = identity // TODO: check enrollment for video_id if provided

	var req clearKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	resp := clearKeyResponse{Keys: []clearKeyKey{}}
	for _, kidB64 := range req.Kids {
		kidHex, err := base64URLToHex(kidB64)
		if err != nil {
			continue
		}
		keyHex, ok := getDashKey(kidHex)
		if !ok {
			continue
		}
		keyB64, err := hexToBase64URL(keyHex)
		if err != nil {
			continue
		}
		resp.Keys = append(resp.Keys, clearKeyKey{
			Kty: "oct",
			Kid: kidB64, // echo back as received
			K:   keyB64,
		})
	}

	if len(resp.Keys) == 0 {
		http.Error(w, "no keys found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// generateDashKey creates a random 128-bit key pair (hex encoded)
func generateDashKey() (keyID, key string, err error) {
	kid := make([]byte, 16)
	if _, err = rand.Read(kid); err != nil {
		return "", "", err
	}
	k := make([]byte, 16)
	if _, err = rand.Read(k); err != nil {
		return "", "", err
	}
	return hex.EncodeToString(kid), hex.EncodeToString(k), nil
}
