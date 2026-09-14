package main

import (
	"go.uber.org/fx"

	"github.com/telemetry-platform/geo-service/internal/core/health"
	"github.com/telemetry-platform/geo-service/internal/core/position"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/api/router"
	httpclient "github.com/telemetry-platform/geo-service/internal/infrastructure/http"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/postgres"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/rabbitmq"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/redis"
)

// Module aggregates all FX modules of the application.
func Module() fx.Option {
	return fx.Options(
		logger.Module(),
		env.Module(),
		postgres.Module(),
		redis.Module(),
		rabbitmq.Module(),
		httpclient.Module(),
		health.Module,
		position.Module,
		router.Module(),
	)
}
