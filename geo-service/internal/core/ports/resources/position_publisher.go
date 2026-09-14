package resources

import (
	"context"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
)

// PositionPublisher is the port for publishing positions to a message queue.
//
//go:generate mockery --name PositionPublisher --dir=. --output=./mocks
type PositionPublisher interface {
	Publish(ctx context.Context, pos domain.Position) error
}
