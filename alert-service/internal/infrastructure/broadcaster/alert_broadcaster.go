package broadcaster

import (
	"context"
	"sync"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/core/ports/resources"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// subscriberBufferSize is the buffer size for subscriber channels.
const subscriberBufferSize = 16

// Compile-time check that Broadcaster implements AlertBroadcaster.
var _ resources.AlertBroadcaster = (*Broadcaster)(nil)

// Broadcaster implements resources.AlertBroadcaster with an in-memory
// list of subscriber channels protected by a mutex.
type Broadcaster struct {
	mu          sync.Mutex
	subscribers []chan domain.Alert
	logger      logger.Logger
}

// NewBroadcaster creates and returns a new in-memory Broadcaster.
func NewBroadcaster(log logger.Logger) *Broadcaster {
	return &Broadcaster{logger: log}
}

// Broadcast sends the alert to all subscribers in a non-blocking manner.
// If a subscriber's channel is full, the alert is dropped for that subscriber
// and a warning is logged.
func (b *Broadcaster) Broadcast(ctx context.Context, alert domain.Alert) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ch := range b.subscribers {
		select {
		case ch <- alert:
		default:
			b.logger.Warnw(ctx, "dropping alert for slow subscriber",
				"vehicle_id", alert.VehicleID, "type", alert.Type)
		}
	}

	return nil
}

// Subscribe creates a new buffered subscriber channel and adds it to the
// broadcaster. The returned channel receives alerts as they are broadcast.
func (b *Broadcaster) Subscribe() <-chan domain.Alert {
	ch := make(chan domain.Alert, subscriberBufferSize)

	b.mu.Lock()
	b.subscribers = append(b.subscribers, ch)
	b.mu.Unlock()

	return ch
}

// Unsubscribe removes the given subscriber channel from the broadcaster
// and closes it. If the channel is not found, this is a no-op.
func (b *Broadcaster) Unsubscribe(ch <-chan domain.Alert) {
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
