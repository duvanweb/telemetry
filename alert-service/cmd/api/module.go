package main

import (
	"context"

	"go.uber.org/fx"

	"github.com/telemetry-platform/alert-service/internal/core/alert"
	"github.com/telemetry-platform/alert-service/internal/core/health"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/router"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/broadcaster"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/postgres"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/rabbitmq"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/redis"
)

// Module aggregates all FX modules of the application.
func Module() fx.Option {
	return fx.Options(
		logger.Module(),
		env.Module(),
		health.Module,
		alert.Module,
		postgres.Module(),
		redis.Module(),
		rabbitmq.Module(),
		broadcaster.Module(),
		router.Module(),
		fx.Invoke(registerConsumerHooks),
		fx.Invoke(registerVehicleDeletionConsumerHooks),
	)
}

// registerConsumerHooks starts and stops the RabbitMQ consumer via FX lifecycle.
func registerConsumerHooks(lc fx.Lifecycle, consumer *rabbitmq.Consumer, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infow(ctx, "starting alert positions consumer")
			return consumer.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			log.Infow(ctx, "stopping alert positions consumer")
			return consumer.Stop()
		},
	})
}

// registerVehicleDeletionConsumerHooks starts and stops the vehicle deletion consumer via FX lifecycle.
func registerVehicleDeletionConsumerHooks(lc fx.Lifecycle, consumer *rabbitmq.VehicleDeletionConsumer, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infow(ctx, "starting vehicle deletion consumer")
			return consumer.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			log.Infow(ctx, "stopping vehicle deletion consumer")
			return consumer.Stop()
		},
	})
}
