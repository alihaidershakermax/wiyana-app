package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/bot"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/cache"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/config"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/db"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/security"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/store"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/video"
)

type server struct {
	cfg      config.Config
	data     store.Store
	catalog  video.Catalog
	adb      *db.DB
	security *security.Service
	streamer *video.Streamer
	spa      http.Handler
	limiter  *rateLimiter
	cache    cache.Cache
}

type rateLimiter struct {
	mu   sync.Mutex
	hits map[string]windowHit
}
type windowHit struct {
	started time.Time
	count   int
}

func newRateLimiter() *rateLimiter { return &rateLimiter{hits: map[string]windowHit{}} }

func (l *rateLimiter) allow(key string, limit int, window time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	hit := l.hits[key]
	if hit.started.IsZero() || now.Sub(hit.started) >= window {
		l.hits[key] = windowHit{started: now, count: 1}
		return true
	}
	if hit.count >= limit {
		return false
	}
	hit.count++
	l.hits[key] = hit
	return true
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	data := createStore(cfg)
	defer data.Close()
	responseCache := createCache(cfg)
	defer responseCache.Close()
	catalog := createCatalog(cfg)
	// Keep a reference to the underlying Convex catalog before caching
	// wraps it — the Telegram bot needs the generic Convex call surface
	// (users, share links, broadcasts) that only the Convex catalog exposes.
	// When SQLite is the catalog, use the SQLite-backed bot store instead.
	var botStore bot.ConvexStore
	if cv, ok := catalog.(bot.ConvexStore); ok {
		botStore = cv
	} else if sc := video.UnwrapSQLite(catalog); sc != nil {
		botStore = bot.NewSQLiteStore(sc.DB())
	}
	catalog = video.NewCachedCatalog(catalog, responseCache, time.Duration(cfg.CacheTTLSeconds*2)*time.Second, 5*time.Second)
	defer catalog.Close()
	s := &server{cfg: cfg, data: data, catalog: catalog, cache: responseCache, limiter: newRateLimiter()}
	if sc := video.UnwrapSQLite(catalog); sc != nil {
		s.adb = sc.DB()
	}
	s.security = security.New(cfg, data, catalog)
	s.streamer = video.NewStreamer(cfg.FFmpeg)
	if info, err := os.Stat("static/spa"); err == nil && info.IsDir() {
		s.spa = newSPAServer("static/spa")
	}
	h := s.middleware(s)
	address := net.JoinHostPort(cfg.Host, cfg.Port)
	serverHTTP := &http.Server{Addr: address, Handler: h, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 60 * time.Second, IdleTimeout: 120 * time.Second}
	log.Printf("Go API listening on %s", address)

	// Telegram bot runs alongside the HTTP server when a bot token is
	// configured. It uses the same Convex catalog (polling, no HTTP port needed).
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if cfg.BotsEnabled && cfg.BotToken != "" {
		cfg := bot.Config{
			Token:     cfg.BotToken,
			WebAppURL: cfg.WebAppURL,
			OwnerID:   config.BotOwnerID(),
		}
		// Only the Convex catalog exposes the generic Convex call surface the
		// bot needs (users, share links, videos). Fall back gracefully when
		// a non-Convex catalog is in use.
		if botStore != nil {
			b := bot.New(cfg, botStore)
			go func() {
				if err := b.Run(ctx); err != nil {
					log.Printf("WARNING: telegram bot stopped (%v); the API remains available without the bot", err)
				}
			}()
			log.Printf("telegram bot starting (polling mode)")
		} else {
			log.Printf("WARNING: telegram bot skipped; current catalog (%T) does not support bot operations. Use CONVEX_URL as the catalog backend to enable the bot.", catalog)
		}
	}

	if err := serverHTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

// redisRetryPolicy returns how many times and how long to wait when the Redis
// dial fails at boot. In production without an explicit REDIS_URL the default
// points at localhost:6379 which never exists; failing fast repeatedly wastes
// restarts, so the service waits and retries with backoff before giving up.
func redisRetryPolicy(inProduction bool) (maxAttempts int, wait time.Duration) {
	if inProduction {
		// Redis is now optional: after exhausting retries the service falls
		// back to an in-memory store and keeps serving. Keep the wait short
		// so a broken REDIS_URL never stalls boot for more than ~25s.
		return 3, 10 * time.Second // total wait ~25s: 0,10,20
	}
	return 2, 5 * time.Second
}

func createStore(cfg config.Config) store.Store {
	if cfg.StoreBackend == "memory" {
		log.Printf("WARNING: memory store enabled; not suitable for multiple production instances")
		return store.NewMemory()
	}
	maxAttempts, wait := redisRetryPolicy(strings.EqualFold(cfg.Environment, "production"))
	inProduction := strings.EqualFold(cfg.Environment, "production")
	var data store.Store
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		data, err = store.NewRedis(cfg.RedisURL)
		if err == nil {
			return data
		}
		if attempt == 1 {
			log.Printf("Redis dial failed: %v", err)
			if !strings.Contains(cfg.RedisURL, "://localhost") && !strings.Contains(cfg.RedisURL, "://127.0.0.1") {
				log.Printf("hint: verify REDIS_URL=%s is reachable from this environment", cfg.RedisURL)
			} else {
				log.Printf("hint: REDIS_URL is unset; defaulting to localhost:6379. Add REDIS_URL (e.g. Redis on Render / Upstash) to fix this")
			}
		}
		if attempt < maxAttempts {
			log.Printf("retrying Redis connection in %s (attempt %d/%d)...", wait, attempt+1, maxAttempts)
			time.Sleep(wait)
		}
	}
	if inProduction {
		log.Printf("WARNING: Redis unavailable after %d attempts (%v); falling back to in-memory store. Persisted data will NOT be shared across instances; fix REDIS_URL or create a Redis service to restore Redis storage.", maxAttempts, err)
	} else {
		log.Printf("Redis unavailable after %d attempts; using memory store for development: %v", maxAttempts, err)
	}
	return store.NewMemory()
}

func createCache(cfg config.Config) cache.Cache {
	if !cfg.CacheEnabled {
		log.Printf("cache disabled by configuration")
		return cache.NewMemory(cfg.CacheMaxEntries)
	}
	if cfg.StoreBackend != "memory" {
		responseCache, err := cache.NewRedis(cfg.RedisURL, cfg.CachePrefix, cfg.RedisPoolSize)
		if err == nil {
			return responseCache
		}
		if strings.EqualFold(cfg.Environment, "production") {
			log.Printf("WARNING: Redis cache unavailable (%v); using bounded in-memory response cache: %v", cfg.RedisURL, err)
		} else {
			log.Printf("Redis cache unavailable; using bounded memory cache for development: %v", err)
		}
	}
	return cache.NewMemory(cfg.CacheMaxEntries)
}

func createCatalog(cfg config.Config) video.Catalog {
	// SQLite is the primary database (local file, zero-config).
	if catalog, err := video.NewSQLiteCatalog(cfg.SQLitePath); err != nil {
		log.Printf("WARNING: sqlite catalog unavailable (%v); trying legacy backends.", err)
	} else {
		log.Printf("catalog: using SQLite at %s", cfg.SQLitePath)
		return catalog
	}
	if cfg.ConvexURL != "" {
		catalog := video.NewConvex(cfg.ConvexURL, cfg.ConvexAdminKey, cfg.ConvexBackendKey)
		if err := catalog.Ping(context.Background()); err != nil {
			if strings.EqualFold(cfg.Environment, "production") {
				log.Printf("WARNING: convex catalog unavailable (%v); starting with an empty in-memory catalog. Add a valid CONVEX_URL/CONVEX_ADMIN_KEY to restore the full video catalog.", err)
			} else {
				log.Printf("convex catalog unavailable; using empty memory catalog: %v", err)
			}
		} else {
			return catalog
		}
	}
	if cfg.DatabaseURL != "" {
		catalog, err := video.NewSQL(cfg.DatabaseURL, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, time.Duration(cfg.DBConnMaxLifetimeSecs)*time.Second)
		if err != nil {
			if strings.EqualFold(cfg.Environment, "production") {
				log.Printf("WARNING: database unavailable (%v); starting with an empty in-memory catalog. Add a valid DATABASE_URL to restore the full video catalog.", err)
			} else {
				log.Printf("database unavailable; using empty memory catalog: %v", err)
			}
		} else {
			return catalog
		}
	}
	if strings.EqualFold(cfg.Environment, "production") {
		log.Printf("WARNING: no database configured; running production with an empty in-memory catalog. Videos will not be persisted across restarts.")
	}
	return video.NewMemory(cfg.UploadDir)
}
