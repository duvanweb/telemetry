package cache

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/telemetry-platform/geo-service/internal/core/ports/resources"
)

// Compile-time check that Cache implements PositionCache.
var _ resources.PositionCache = (*Cache)(nil)

// Cache implements resources.PositionCache with Redis.
type Cache struct {
	client *goredis.Client
}

// NewCache creates and returns a new position Cache.
func NewCache(client *goredis.Client) *Cache {
	return &Cache{client: client}
}

// Exists checks if a key exists in Redis.
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}
	return n > 0, nil
}

// Set sets a key in Redis with the given value and TTL.
func (c *Cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}
	return nil
}
