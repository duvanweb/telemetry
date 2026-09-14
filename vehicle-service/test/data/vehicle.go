package testdata

import (
	"time"

	"github.com/telemetry-platform/vehicle-service/internal/core/domain"
)

// GetTestVehicle returns a test vehicle with typical values.
func GetTestVehicle() *domain.Vehicle {
	return &domain.Vehicle{
		ID:        1,
		Plate:     "ABC-123",
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

// GetTestVehicles returns a list of test vehicles.
func GetTestVehicles() []domain.Vehicle {
	return []domain.Vehicle{
		{
			ID:        1,
			Plate:     "ABC-123",
			CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        2,
			Plate:     "DEF-456",
			CreatedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}
}

// GetTestVehicleID returns a test vehicle ID.
func GetTestVehicleID() int64 {
	return 1
}
