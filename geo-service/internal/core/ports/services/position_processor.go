package services

import (
	"context"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
)

// PositionProcessor is the interface for processing GPS positions from the queue.
// It wraps persistence with a circuit breaker to handle transient database failures.

//go:generate mockery --name PositionProcessor --dir=. --output=./mocks
type PositionProcessor interface {
	Process(ctx context.Context, pos domain.Position) error
}
