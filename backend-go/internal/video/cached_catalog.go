package video

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/cache"
)

type cachedCatalog struct {
	base      Catalog
	cache     cache.Cache
	itemTTL   time.Duration
	accessTTL time.Duration
}

func NewCachedCatalog(base Catalog, responseCache cache.Cache, itemTTL, accessTTL time.Duration) Catalog {
	if base == nil || responseCache == nil {
		return base
	}
	return &cachedCatalog{base: base, cache: responseCache, itemTTL: itemTTL, accessTTL: accessTTL}
}

// Inner returns the wrapped catalog (used to reach the SQLite store).
func (c *cachedCatalog) Inner() Catalog { return c.base }

// UnwrapSQLite returns the *SQLiteCatalog behind c, unwrapping caches.
func UnwrapSQLite(c Catalog) *SQLiteCatalog {
	for c != nil {
		if sc, ok := c.(*SQLiteCatalog); ok {
			return sc
		}
		if w, ok := c.(interface{ Inner() Catalog }); ok {
			c = w.Inner()
			continue
		}
		return nil
	}
	return nil
}

func (c *cachedCatalog) Close() error { return c.base.Close() }

func (c *cachedCatalog) Find(ctx context.Context, id string) (Item, error) {
	key := "catalog:item:" + id
	if raw, ok, err := c.cache.Get(ctx, key); err == nil && ok {
		var item Item
		if json.Unmarshal(raw, &item) == nil {
			return item, nil
		}
	}
	item, err := c.base.Find(ctx, id)
	if err != nil {
		return Item{}, err
	}
	if raw, err := json.Marshal(item); err == nil {
		_ = c.cache.Set(ctx, key, raw, c.itemTTL)
	}
	return item, nil
}

func (c *cachedCatalog) CanAccess(ctx context.Context, userID string, item Item) (bool, error) {
	if item.Free {
		return true, nil
	}
	sum := sha256.Sum256([]byte(userID + ":" + item.ID + ":" + item.CourseID))
	key := "catalog:access:" + hex.EncodeToString(sum[:])
	if raw, ok, err := c.cache.Get(ctx, key); err == nil && ok {
		var allowed bool
		if json.Unmarshal(raw, &allowed) == nil {
			return allowed, nil
		}
	}
	allowed, err := c.base.CanAccess(ctx, userID, item)
	if err != nil {
		return false, err
	}
	if raw, err := json.Marshal(allowed); err == nil {
		_ = c.cache.Set(ctx, key, raw, c.accessTTL)
	}
	return allowed, nil
}
