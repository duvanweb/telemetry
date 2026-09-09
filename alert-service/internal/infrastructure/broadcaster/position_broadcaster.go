package broadcaster

import (
	"context"
	"sync"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/core/ports/resources"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// Compile-time check that PositionBroadcaster implements resources.PositionBroadcaster.
var _ resources.PositionBroadcaster = (*PositionBroadcaster)(nil)

// PositionBroadcaster implements resources.PositionBroadcaster with an in-memory
// list of subscriber channels protected by a mutex.
type PositionBroadcaster struct {
	mu          sync.Mutex
	subscribers []chan domain.Position
	logger      logger.Logger
}

// NewPositionBroadcaster creates and returns a new in-memory PositionBroadcaster.
func NewPositionBroadcaster(log logger.Logger) *PositionBroadcaster {
	return &PositionBroadcaster{logger: log}
}

// Broadcast sends the position to all subscribers in a non-blocking manner.
// If a subscriber's channel is full, the position is dropped for that subscriber
// and a warning is logged.
func (b *PositionBroadcaster) Broadcast(ctx context.Context, pos domain.Position) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ch := range b.subscribers {
		select {
		case ch <- pos:
		default:
			b.logger.Warnw(ctx, "dropping position for slow subscriber",
				"vehicle_id", pos.VehicleID)
		}
	}

	return nil
}

// Subscribe creates a new buffered subscriber channel and adds it to the
// broadcaster. The returned channel receives positions as they are broadcast.
func (b *PositionBroadcaster) Subscribe() <-chan domain.Position {
	ch := make(chan domain.Position, subscriberBufferSize)

	b.mu.Lock()
	b.subscribers = append(b.subscribers, ch)
	b.mu.Unlock()

	return ch
}

// Unsubscribe removes the given subscriber channel from the broadcaster
// and closes it. If the channel is not found, this is a no-op.
func (b *PositionBroadcaster) Unsubscribe(ch <-chan domain.Position) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, sub := range b.subscribers {
		if sub == ch {
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			close(sub)
			return
		}
	}
}
