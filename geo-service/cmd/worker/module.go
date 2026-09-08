package main

import (
	"context"

	"go.uber.org/fx"

	"github.com/telemetry-platform/geo-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/postgres"
	positionrepo "github.com/telemetry-platform/geo-service/internal/infrastructure/postgres/repositories/position"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/rabbitmq"
)

// Module aggregates all FX modules for the worker binary.
func Module() fx.Option {
	return fx.Options(
		logger.Module(),
		env.Module(),
		postgres.Module(),
		rabbitmq.Module(),
		fx.Provide(
			fx.Annotate(positionrepo.NewRepository, fx.As(new(repositories.PositionRepository))),
		),
		fx.Provide(rabbitmq.NewConsumer),
		fx.Invoke(registerConsumerHooks),
	)
}

// registerConsumerHooks starts and stops the RabbitMQ consumer via FX lifecycle.
func registerConsumerHooks(lc fx.Lifecycle, consumer *rabbitmq.Consumer, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infow(ctx, "starting position consumer")
			return consumer.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			log.Infow(ctx, "stopping position consumer")
			return consumer.Stop()
		},
	})
}
