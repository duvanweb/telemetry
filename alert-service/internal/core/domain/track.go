package domain

import "time"

// VehicleTrack holds the tracking state for a vehicle used by the anomaly
// detection service to detect the "Vehicle Stopped" anomaly.
type VehicleTrack struct {
	VehicleID   int64
	Latitude    float64
	Longitude   float64
	FirstSeenAt time.Time
	Alerted     bool
}
