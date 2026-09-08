package rabbitmq

import (
	"context"

	"go.uber.org/fx"

	"github.com/telemetry-platform/geo-service/internal/core/ports/resources"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// Module provides the RabbitMQ client and publisher via FX.
func Module() fx.Option {
	return fx.Module(
		"rabbitmq",
		fx.Provide(
			NewClient,
			fx.Annotate(NewPublisher, fx.As(new(resources.PositionPublisher))),
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
