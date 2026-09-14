package domain

import "time"

// AlertType represents the type of an alert.
type AlertType string

// AlertTypeVehicleStopped indicates a vehicle has been stopped at the same
// coordinates for longer than the configured threshold.
const AlertTypeVehicleStopped AlertType = "VEHICLE_STOPPED"

// Alert represents an anomaly alert detected by the alert-service.
type Alert struct {
	ID         int64
	VehicleID  int64
	Type       AlertType
	Latitude   float64
	Longitude  float64
	DetectedAt time.Time
	CreatedAt  time.Time
}
