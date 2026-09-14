package processor

import (
	"context"
	"errors"
	"time"

	"github.com/sony/gobreaker"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	"github.com/telemetry-platform/geo-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// Processor implements services.PositionProcessor by persisting positions
// through a circuit breaker to handle transient database failures.
type Processor struct {
	repo   repositories.PositionRepository
	cb     *gobreaker.CircuitBreaker
	logger logger.Logger
}

// NewProcessor creates a new Processor with a circuit breaker.
// The circuit breaker opens after failureThreshold consecutive failures
// and transitions to half-open after timeout.
func NewProcessor(
	repo repositories.PositionRepository,
	failureThreshold uint32,
	timeout time.Duration,
	log logger.Logger,
) *Processor {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name: "position-persistence",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= failureThreshold
		},
		Timeout: timeout,
	})

	return &Processor{
		repo:   repo,
		cb:     cb,
		logger: log,
	}
}

// Process persists a position through the circuit breaker.
// Returns domain.ErrCircuitBreakerOpen when the circuit breaker is open
// and the repository call is skipped.
func (p *Processor) Process(ctx context.Context, pos domain.Position) error {
	_, err := p.cb.Execute(func() (any, error) {
		return nil, p.repo.Save(ctx, pos)
	})
	if errors.Is(err, gobreaker.ErrOpenState) {
		p.logger.Warnw(ctx, "circuit breaker is open, skipping persistence", "vehicle_id", pos.VehicleID)
		return domain.ErrCircuitBreakerOpen
	}
	return err
}
