package testdata

import (
	"time"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
)

// GetTestVehicleTrack returns a test vehicle track that has not been alerted.
func GetTestVehicleTrack() domain.VehicleTrack {
	return domain.VehicleTrack{
		VehicleID:   1,
		Latitude:    4.71,
		Longitude:   -74.07,
		FirstSeenAt: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
		Alerted:     false,
	}
}

// GetTestAlertedTrack returns a test vehicle track that has already been alerted.
func GetTestAlertedTrack() domain.VehicleTrack {
	track := GetTestVehicleTrack()
	track.Alerted = true
	return track
}
