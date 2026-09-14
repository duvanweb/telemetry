package testdata

import (
	"time"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
)

// GetTestPosition returns a test position with typical values.
func GetTestPosition() domain.Position {
	return domain.Position{
		VehicleID:  1,
		Latitude:   4.71,
		Longitude:  -74.07,
		RecordedAt: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}

// GetTestInvalidPosition returns a test position with invalid latitude.
func GetTestInvalidPosition() domain.Position {
	return domain.Position{
		VehicleID:  1,
		Latitude:   100.0,
		Longitude:  -74.07,
		RecordedAt: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}
