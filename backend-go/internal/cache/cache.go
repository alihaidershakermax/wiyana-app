package cache

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/netutil"
	"github.com/redis/go-redis/v9"
)

var ErrUnavailable = errors.New("cache unavailable")

type Stats struct {
	Hits   uint64 `json:"hits"`
	Misses uint64 `json:"misses"`
	Items  int    `json:"items"`
}

type Cache interface {
	Get(context.Context, string) ([]byte, bool, error)
	Set(context.Context, string, []byte, time.Duration) error
	Delete(context.Context, string) error
	DeletePrefix(context.Context, string) error
	Stats() Stats
	Close() error
}

type RedisCache struct {
	client *redis.Client
	prefix string
	hits   atomic.Uint64
	misses atomic.Uint64
}

func NewRedis(redisURL, prefix string, poolSize int) (*RedisCache, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	if poolSize > 0 {
		options.PoolSize = poolSize
	}
	// IPv4-only connection: Render (Oregon) cannot reach external IPv6
	// addresses; DNS may return AAAA records first.
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
	return &RedisCache{client: client, prefix: prefix}, nil
}

func (c *RedisCache) key(key string) string { return c.prefix + key }

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	value, err := c.client.Get(ctx, c.key(key)).Bytes()
	if errors.Is(err, redis.Nil) {
		c.misses.Add(1)
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	c.hits.Add(1)
	return value, true, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl <= 0 {
		return c.client.Set(ctx, c.key(key), value, 0).Err()
	}
	return c.client.Set(ctx, c.key(key), value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.key(key)).Err()
}

func (c *RedisCache) DeletePrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	pattern := c.key(prefix) + "*"
	for {
		keys, next, err := c.client.Scan(ctx, cursor, pattern, 256).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}

func (c *RedisCache) Stats() Stats {
	return Stats{Hits: c.hits.Load(), Misses: c.misses.Load()}
}

func (c *RedisCache) Close() error { return c.client.Close() }

type memoryEntry struct {
	value     []byte
	expiresAt time.Time
	usedAt    uint64
}

type MemoryCache struct {
	mu         sync.Mutex
	items      map[string]memoryEntry
	maxEntries int
	clock      uint64
	hits       uint64
	misses     uint64
}

func NewMemory(maxEntries int) *MemoryCache {
	if maxEntries < 64 {
		maxEntries = 1024
	}
	return &MemoryCache{items: make(map[string]memoryEntry), maxEntries: maxEntries}
}

func (c *MemoryCache) Get(_ context.Context, key string) ([]byte, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.items[key]
	if !ok || (!entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt)) {
		if ok {
			delete(c.items, key)
		}
		c.misses++
		return nil, false, nil
	}
	c.clock++
	entry.usedAt = c.clock
	c.items[key] = entry
	c.hits++
	return append([]byte(nil), entry.value...), true, nil
}

func (c *MemoryCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clock++
	if len(c.items) >= c.maxEntries {
		c.evictOne()
	}
	entry := memoryEntry{value: append([]byte(nil), value...), usedAt: c.clock}
	if ttl > 0 {
		entry.expiresAt = time.Now().Add(ttl)
	}
	c.items[key] = entry
	return nil
}

func (c *MemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *MemoryCache) DeletePrefix(_ context.Context, prefix string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key := range c.items {
		if strings.HasPrefix(key, prefix) {
			delete(c.items, key)
		}
	}
	return nil
}

func (c *MemoryCache) Stats() Stats {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Stats{Hits: c.hits, Misses: c.misses, Items: len(c.items)}
}

func (c *MemoryCache) Close() error { return nil }

func (c *MemoryCache) evictOne() {
	var oldestKey string
	var oldest uint64
	for key, entry := range c.items {
		if oldestKey == "" || entry.usedAt < oldest {
			oldestKey = key
			oldest = entry.usedAt
		}
	}
	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}
