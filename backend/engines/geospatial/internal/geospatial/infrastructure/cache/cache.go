package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheEntry
}

type cacheEntry struct {
	Value    []byte
	ExpireAt time.Time
}

func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{items: make(map[string]cacheEntry)}
	go c.cleanupLoop()
	return c
}

func (c *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.items[key]
	if !ok || time.Now().After(entry.ExpireAt) {
		return nil, fmt.Errorf("cache miss")
	}
	return entry.Value, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheEntry{Value: value, ExpireAt: time.Now().Add(ttl)}
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.items[key]
	return ok, nil
}

func (c *MemoryCache) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.items {
			if now.After(v.ExpireAt) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

type GeospatialCache struct {
	inner MemoryCache
	ttl   time.Duration
}

func NewGeospatialCache(ttl time.Duration) *GeospatialCache {
	return &GeospatialCache{inner: *NewMemoryCache(), ttl: ttl}
}

func (gc *GeospatialCache) GetOrFetch(ctx context.Context, key string, fetch func() (any, error)) (any, error) {
	data, err := gc.inner.Get(ctx, key)
	if err == nil && data != nil {
		var result any
		json.Unmarshal(data, &result)
		if result != nil {
			return result, nil
		}
	}

	result, err := fetch()
	if err != nil {
		return nil, err
	}

	encoded, _ := json.Marshal(result)
	gc.inner.Set(ctx, key, encoded, gc.ttl)
	return result, nil
}
