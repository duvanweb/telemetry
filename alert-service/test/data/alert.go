package testdata

import (
	"time"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
)

// GetTestAlert returns a test alert with typical values.
func GetTestAlert() domain.Alert {
	return domain.Alert{
		ID:         1,
		VehicleID:  1,
		Type:       domain.AlertTypeVehicleStopped,
		Latitude:   4.71,
		Longitude:  -74.07,
		DetectedAt: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
		CreatedAt:  time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}
