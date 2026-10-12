package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Environment           string
	Host                  string
	Port                  string
	JWTSecret             string
	BotToken              string
	DatabaseURL           string
	SQLitePath            string
	RedisURL              string
	ConvexURL             string
	ConvexAdminKey        string
	ConvexBackendKey      string
	StoreBackend          string
	UploadDir             string
	HLSDir                string
	FFmpeg                string
	WatermarkEnabled      bool
	GrantTTLSeconds       int
	HLSTTLSeconds         int
	SessionTTLSeconds     int
	WarningThreshold      int
	BanThreshold          int
	BanDurationSeconds    int
	BehaviorWindowSeconds int
	EscalationViolations  int
	ProtectionEnabled     bool
	AllowedOrigins        []string
	AllowDevAuth          bool
	CacheEnabled          bool
	CachePrefix           string
	CacheTTLSeconds       int
	CacheMaxEntries       int
	RedisPoolSize         int
	DBMaxOpenConns        int
	DBMaxIdleConns        int
	DBConnMaxLifetimeSecs int
	HTTPRateLimit         int
	BotsEnabled           bool
	WebAppURL             string
	BotAdminIDs           []int64
	GoogleDriveSA         string
	GoogleDriveAPIKey     string
}

func Load() (Config, error) {
	c := Config{
		Environment:           env("ENVIRONMENT", "development"),
		Host:                  env("HOST", "0.0.0.0"),
		Port:                  env("PORT", "8000"),
		JWTSecret:             firstEnv("JWT_SECRET", "SECRET_KEY"),
		BotToken:              firstEnv("TELEGRAM_BOT_TOKEN", "BOT_TOKEN"),
		DatabaseURL:           firstEnv("DATABASE_URL", "SUPABASE_DB_URL", "SUPABASE_POOLER_URL"),
		SQLitePath:            env("SQLITE_PATH", "data/wiyana.db"),
		RedisURL:              env("REDIS_URL", "redis://localhost:6379/0"),
		ConvexURL:             env("CONVEX_URL", ""),
		ConvexAdminKey:        env("CONVEX_ADMIN_KEY", ""),
		ConvexBackendKey:      env("CONVEX_BACKEND_KEY", ""),
		StoreBackend:          strings.ToLower(env("STORE_BACKEND", "redis")),
		UploadDir:             env("UPLOAD_DIR", "uploads"),
		HLSDir:                env("HLS_OUTPUT_DIR", "hls_cache"),
		FFmpeg:                env("FFMPEG_BINARY", "ffmpeg"),
		WatermarkEnabled:      envBool("FFMPEG_WATERMARK_ENABLED", true),
		GrantTTLSeconds:       envInt("PROTECTION_GRANT_TTL_SECONDS", 300),
		HLSTTLSeconds:         envInt("PROTECTION_HLS_TTL_SECONDS", 300),
		SessionTTLSeconds:     envInt("PLAYBACK_SESSION_TTL_SECONDS", 300),
		WarningThreshold:      envInt("PROTECTION_WARNING_THRESHOLD", 5),
		BanThreshold:          envInt("PROTECTION_BAN_THRESHOLD", 10),
		BanDurationSeconds:    envInt("PROTECTION_BAN_DURATION_SECONDS", 30*24*60*60),
		BehaviorWindowSeconds: envInt("PROTECTION_BEHAVIORAL_WINDOW_SECONDS", 300),
		EscalationViolations:  envInt("PROTECTION_ESCALATION_VIOLATIONS", 3),
		ProtectionEnabled:     envBool("PROTECTION_ENABLED", true),
		AllowedOrigins:        splitCSV(env("CORS_ORIGINS", "")),
		AllowDevAuth:          envBool("ALLOW_DEV_AUTH", false),
		CacheEnabled:          envBool("CACHE_ENABLED", true),
		CachePrefix:           env("CACHE_PREFIX", "wiyana:cache:"),
		CacheTTLSeconds:       envInt("CACHE_TTL_SECONDS", 30),
		CacheMaxEntries:       envInt("CACHE_MAX_ENTRIES", 4096),
		RedisPoolSize:         envInt("REDIS_POOL_SIZE", 32),
		DBMaxOpenConns:        envInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:        envInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetimeSecs: envInt("DB_CONN_MAX_LIFETIME_SECONDS", 1800),
		HTTPRateLimit:         envInt("HTTP_RATE_LIMIT_PER_MINUTE", 120),
		BotsEnabled:           envBool("BOTS_ENABLED", true),
		WebAppURL:             env("WEBAPP_URL", "https://wena-app.onrender.com/"),
		GoogleDriveSA:         env("GOOGLE_DRIVE_SA_JSON", ""),
		GoogleDriveAPIKey:     env("GOOGLE_DRIVE_API_KEY", ""),
		BotAdminIDs:           splitInts(env("BOT_ADMIN_IDS", fmt.Sprint(BotOwnerID()))),
	}
	if len(c.AllowedOrigins) == 0 {
		c.AllowedOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	}
	if c.JWTSecret == "" {
		generated, err := randomSecret()
		if err != nil {
			return Config{}, fmt.Errorf("generate JWT secret: %w", err)
		}
		if strings.EqualFold(c.Environment, "production") {
			// Production safety net: Render/Docker do not persist secrets that are
			// not declared in the deployment manifest. Generating a cryptographically
			// random secret at boot lets the service start instead of crash-looping,
			// but every restart produces a NEW secret — invalidating all existing
			// sessions and grants. Set JWT_SECRET explicitly when stable sessions
			// across restarts are required.
			log.Printf("WARNING: JWT_SECRET missing in production; generated a random secret at boot. Sessions will be invalidated on every restart. Set JWT_SECRET explicitly to fix this.")
		} else {
			log.Printf("development JWT secret generated (per-restart)")
		}
		c.JWTSecret = generated
	}
	if c.BanThreshold < c.WarningThreshold {
		return Config{}, errors.New("PROTECTION_BAN_THRESHOLD must be >= PROTECTION_WARNING_THRESHOLD")
	}
	if c.CacheTTLSeconds < 1 || c.CacheMaxEntries < 64 {
		return Config{}, errors.New("cache TTL and size values are too small")
	}
	if c.DBMaxIdleConns > c.DBMaxOpenConns {
		return Config{}, errors.New("DB_MAX_IDLE_CONNS must be <= DB_MAX_OPEN_CONNS")
	}
	if c.GrantTTLSeconds < 30 || c.SessionTTLSeconds < 30 {
		return Config{}, errors.New("playback TTL values must be at least 30 seconds")
	}
	return c, nil
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(os.Getenv(key)))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func envBool(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func splitCSV(value string) []string {
	var out []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func randomSecret() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// botOwnerID is the platform owner telegram id (can be overridden via env).
func BotOwnerID() int64 {
	if v := strings.TrimSpace(os.Getenv("BOT_OWNER_ID")); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return n
		}
	}
	return 2018954602 // ali_haider_shaker on Telegram
}

func splitInts(value string) []int64 {
	var out []int64
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		n, err := strconv.ParseInt(item, 10, 64)
		if err == nil {
			out = append(out, n)
		}
	}
	return out
}
