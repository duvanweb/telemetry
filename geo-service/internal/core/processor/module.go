package processor

import (
	"fmt"
	"time"

	"go.uber.org/fx"

	"github.com/telemetry-platform/geo-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/geo-service/internal/core/ports/services"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// Module wires the position processor with circuit breaker into FX.
func Module() fx.Option {
	return fx.Module(
		"processor",
		fx.Provide(
			fx.Annotate(newProcessorFromConfig, fx.As(new(services.PositionProcessor))),
		),
	)
}

// newProcessorFromConfig creates a PositionProcessor with circuit breaker config from env.
func newProcessorFromConfig(repo repositories.PositionRepository, config *env.Configuration, log logger.Logger) (*Processor, error) {
	timeout, err := time.ParseDuration(config.CBTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to parse circuit breaker timeout: %w", err)
	}
	return NewProcessor(repo, uint32(config.CBFailureThreshold), timeout, log), nil
}
