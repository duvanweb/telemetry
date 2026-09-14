package resources

import (
	"context"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
)

// PositionTracker is the interface for tracking vehicle positions in Redis.
//
//go:generate mockery --name PositionTracker --dir=. --output=./mocks
type PositionTracker interface {
	Get(ctx context.Context, vehicleID int64) (domain.VehicleTrack, error)
	Set(ctx context.Context, track domain.VehicleTrack) error
	Delete(ctx context.Context, vehicleID int64) error
}
