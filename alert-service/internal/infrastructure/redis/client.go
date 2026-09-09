package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"

	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// NewRedisClient creates and returns a new Redis client.
// It pings Redis to verify connectivity before returning.
func NewRedisClient(config *env.Configuration, log logger.Logger) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr: config.RedisAddr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Infow(context.Background(), "redis connection established", "addr", config.RedisAddr)

	return client, nil
}
