package rabbitmq

import (
	"context"

	"go.uber.org/fx"

	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// Module provides the RabbitMQ client and consumer via FX.
func Module() fx.Option {
	return fx.Module(
		"rabbitmq",
		fx.Provide(
			NewClient,
			NewConsumer,
		),
		fx.Invoke(registerRabbitMQHooks),
	)
}

// registerRabbitMQHooks registers the RabbitMQ close hook on the FX lifecycle.
func registerRabbitMQHooks(lc fx.Lifecycle, client *Client, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Infow(ctx, "closing rabbitmq connection")
			return client.Close()
		},
	})
}
