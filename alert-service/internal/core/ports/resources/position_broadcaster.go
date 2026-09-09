package resources

import (
	"context"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
)

// PositionBroadcaster is the interface for broadcasting positions to SSE subscribers.
//
//go:generate mockery --name PositionBroadcaster --dir=. --output=./mocks
type PositionBroadcaster interface {
	Broadcast(ctx context.Context, pos domain.Position) error
	Subscribe() <-chan domain.Position
	Unsubscribe(ch <-chan domain.Position)
}
