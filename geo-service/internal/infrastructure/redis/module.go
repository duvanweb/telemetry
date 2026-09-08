package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// Module provides the Redis client via FX.
func Module() fx.Option {
	return fx.Module(
		"redis",
		fx.Provide(NewRedisClient),
		fx.Invoke(registerRedisHooks),
	)
}

// registerRedisHooks registers the Redis close hook on the FX lifecycle.
func registerRedisHooks(lc fx.Lifecycle, client *goredis.Client, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Infow(ctx, "closing redis connection")
			return client.Close()
		},
	})
}
