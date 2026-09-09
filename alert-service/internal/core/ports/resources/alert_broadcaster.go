package resources

import (
	"context"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
)

// AlertBroadcaster is the interface for broadcasting alerts to SSE subscribers.
//
//go:generate mockery --name AlertBroadcaster --dir=. --output=./mocks
type AlertBroadcaster interface {
	Broadcast(ctx context.Context, alert domain.Alert) error
	Subscribe() <-chan domain.Alert
	Unsubscribe(ch <-chan domain.Alert)
}
