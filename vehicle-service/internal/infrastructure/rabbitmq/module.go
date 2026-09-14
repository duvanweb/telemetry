package rabbitmq

import (
	"context"

	"go.uber.org/fx"

	"github.com/telemetry-platform/vehicle-service/internal/core/ports/resources"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
)

// Module wires the RabbitMQ infrastructure into FX.
func Module() fx.Option {
	return fx.Module(
		"rabbitmq",
		fx.Provide(NewClient),
		fx.Provide(
			fx.Annotate(NewPublisher, fx.As(new(resources.EventPublisher))),
		),
		fx.Invoke(registerLifecycleHooks),
	)
}

// registerLifecycleHooks closes the RabbitMQ connection on shutdown.
func registerLifecycleHooks(lc fx.Lifecycle, client *Client, log logger.Logger, config *env.Configuration) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infow(ctx, "rabbitmq client started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Infow(ctx, "closing rabbitmq connection")
			return client.Close()
		},
	})
}
