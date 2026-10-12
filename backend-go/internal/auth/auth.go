package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	JWTSecret   string
	BotToken    string
	Environment string
}

type Identity struct {
	ID         string `json:"id"`
	TelegramID string `json:"telegram_id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
}

var ErrUnauthorized = errors.New("unauthorized")

func Resolve(r *http.Request, cfg Config) (Identity, error) {
	if raw := strings.TrimSpace(r.Header.Get("X-Telegram-Init-Data")); raw != "" {
		if identity, err := VerifyTelegramInitData(raw, cfg.BotToken, cfg.Environment); err == nil {
			return identity, nil
		}
	}
	if header := strings.TrimSpace(r.Header.Get("Authorization")); strings.HasPrefix(strings.ToLower(header), "bearer ") {
		token := strings.TrimSpace(header[len("Bearer "):])
		if identity, err := ParseJWT(token, cfg.JWTSecret); err == nil {
			return identity, nil
		}
	}
	return Identity{}, ErrUnauthorized
}

func ParseJWT(token, secret string) (Identity, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || secret == "" {
		return Identity{}, ErrUnauthorized
	}
	unsigned := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Identity{}, ErrUnauthorized
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(unsigned))
	if !hmac.Equal(signature, mac.Sum(nil)) {
		return Identity{}, ErrUnauthorized
	}
	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	var claims map[string]any
	if err := decodePart(parts[0], &header); err != nil {
		return Identity{}, ErrUnauthorized
	}
	if err := decodePart(parts[1], &claims); err != nil {
		return Identity{}, ErrUnauthorized
	}
	if header.Alg != "HS256" {
		return Identity{}, ErrUnauthorized
	}
	if exp, ok := numberClaim(claims["exp"]); ok && time.Now().Unix() >= exp {
		return Identity{}, ErrUnauthorized
	}
	id := firstClaim(claims, "telegram_id", "user_id", "sub", "id")
	if id == "" {
		return Identity{}, ErrUnauthorized
	}
	return Identity{ID: id, TelegramID: firstClaim(claims, "telegram_id", "sub", "id"), Name: firstClaim(claims, "name", "username", "first_name"), Role: firstClaim(claims, "role")}, nil
}

func VerifyTelegramInitData(raw, botToken, environment string) (Identity, error) {
	if raw == "" || botToken == "" {
		return Identity{}, ErrUnauthorized
	}
	values, err := url.ParseQuery(raw)
	if err != nil {
		return Identity{}, ErrUnauthorized
	}
	received := values.Get("hash")
	values.Del("hash")
	if received == "" {
		return Identity{}, ErrUnauthorized
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, key+"="+values.Get(key))
	}
	secretMac := hmac.New(sha256.New, []byte("WebAppData"))
	_, _ = secretMac.Write([]byte(botToken))
	secret := secretMac.Sum(nil)
	checkMac := hmac.New(sha256.New, secret)
	_, _ = checkMac.Write([]byte(strings.Join(lines, "\n")))
	expected := hex.EncodeToString(checkMac.Sum(nil))
	if !hmac.Equal([]byte(strings.ToLower(received)), []byte(expected)) {
		return Identity{}, ErrUnauthorized
	}
	if authDate := values.Get("auth_date"); authDate != "" {
		ts, err := strconv.ParseInt(authDate, 10, 64)
		if err != nil {
			return Identity{}, ErrUnauthorized
		}
		age := time.Now().Unix() - ts
		if strings.EqualFold(environment, "production") && (age > 86400 || age < -300) {
			return Identity{}, ErrUnauthorized
		}
	}
	var user struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	}
	if payload := values.Get("user"); payload != "" {
		if err := json.Unmarshal([]byte(payload), &user); err != nil || user.ID == 0 {
			return Identity{}, ErrUnauthorized
		}
	} else {
		return Identity{}, ErrUnauthorized
	}
	name := strings.TrimSpace(strings.TrimSpace(user.FirstName) + " " + strings.TrimSpace(user.LastName))
	if name == "" {
		name = user.Username
	}
	id := strconv.FormatInt(user.ID, 10)
	return Identity{ID: id, TelegramID: id, Name: name, Role: "student"}, nil
}

func SignJWT(identity Identity, secret string, ttl time.Duration) (string, error) {
	if secret == "" || identity.ID == "" {
		return "", ErrUnauthorized
	}
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	claims := map[string]any{"sub": identity.ID, "telegram_id": identity.TelegramID, "name": identity.Name, "role": identity.Role, "iat": time.Now().Unix(), "exp": time.Now().Add(ttl).Unix()}
	encode := func(value any) (string, error) {
		data, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return base64.RawURLEncoding.EncodeToString(data), nil
	}
	h, err := encode(header)
	if err != nil {
		return "", err
	}
	c, err := encode(claims)
	if err != nil {
		return "", err
	}
	unsigned := h + "." + c
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(unsigned))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return unsigned + "." + signature, nil
}

func decodePart(part string, target any) error {
	data, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
func numberClaim(value any) (int64, bool) {
	switch v := value.(type) {
	case float64:
		return int64(v), true
	case json.Number:
		n, err := v.Int64()
		return n, err == nil
	case int64:
		return v, true
	default:
		return 0, false
	}
}
func firstClaim(claims map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := claims[key]; ok {
			switch v := value.(type) {
			case string:
				if strings.TrimSpace(v) != "" {
					return v
				}
			case float64:
				return strconv.FormatInt(int64(v), 10)
			case json.Number:
				return v.String()
			}
		}
	}
	return ""
}

func Require(r *http.Request, cfg Config) (Identity, int) {
	identity, err := Resolve(r, cfg)
	if err != nil {
		return Identity{}, http.StatusUnauthorized
	}
	return identity, 0
}

func ErrorMessage(err error) string {
	if errors.Is(err, ErrUnauthorized) {
		return "تعذر التحقق من هوية المستخدم"
	}
	return fmt.Sprintf("%v", err)
}
