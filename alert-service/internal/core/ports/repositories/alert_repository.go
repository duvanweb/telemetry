package repositories

import (
	"context"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
)

// AlertRepository is the interface for persisting and querying alerts.
//
//go:generate mockery --name AlertRepository --dir=. --output=./mocks
type AlertRepository interface {
	Save(ctx context.Context, alert *domain.Alert) error
	List(ctx context.Context, limit, offset int) ([]domain.Alert, int64, error)
	DeleteByVehicleID(ctx context.Context, vehicleID int64) error
}
