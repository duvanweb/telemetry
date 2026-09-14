package resources

import (
	"context"
	"time"
)

// EventPublisher is the port for publishing domain events to the message broker.
//
//go:generate mockery --name EventPublisher --dir=. --output=./mocks
type EventPublisher interface {
	// PublishVehicleDeleted publishes a vehicle.deleted event with the given
	// vehicle ID and deletion timestamp. Consumers use this to clean up
	// orphaned data in their own databases and caches.
	PublishVehicleDeleted(ctx context.Context, vehicleID int64, deletedAt time.Time) error
}
