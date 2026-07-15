package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type GeospatialCache struct {
	inner Cache
	ttl   time.Duration
}

func NewGeospatialCache(inner Cache, ttl time.Duration) *GeospatialCache {
	return &GeospatialCache{inner: inner, ttl: ttl}
}

func (c *GeospatialCache) GetOrFetch(ctx context.Context, key string, fetch func() ([]byte, error)) ([]byte, error) {
	data, err := c.inner.Get(ctx, key)
	if err == nil && data != nil {
		return data, nil
	}

	data, err = fetch()
	if err != nil {
		return nil, err
	}

	c.inner.Set(ctx, key, data, c.ttl)
	return data, nil
}

func (c *GeospatialCache) Invalidate(ctx context.Context, key string) error {
	return c.inner.Delete(ctx, key)
}
