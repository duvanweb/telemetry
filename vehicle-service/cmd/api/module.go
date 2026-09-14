package main

import (
	"go.uber.org/fx"

	"github.com/telemetry-platform/vehicle-service/internal/core/health"
	"github.com/telemetry-platform/vehicle-service/internal/core/vehicle"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/api/router"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/postgres"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/rabbitmq"
)

// Module aggregates all FX modules of the application.
func Module() fx.Option {
	return fx.Options(
		logger.Module(),
		env.Module(),
		postgres.Module(),
		rabbitmq.Module(),
		health.Module,
		vehicle.Module,
		router.Module(),
	)
}
