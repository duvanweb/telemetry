package main

import (
	"go.uber.org/fx"

	"github.com/telemetry-platform/vehicle-service/internal/core/health"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/api/router"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
)

// Module aggregates all FX modules of the application.
func Module() fx.Option {
	return fx.Options(
		logger.Module(),
		env.Module(),
		health.Module,
		router.Module(),
	)
}
