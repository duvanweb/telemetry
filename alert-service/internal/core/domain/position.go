package domain

import "time"

// Position is a GPS position consumed from the alert_positions queue.
// JSON tags match the structure published by geo-service (PascalCase without tags).
type Position struct {
	VehicleID  int64     `json:"VehicleID"`
	Latitude   float64   `json:"Latitude"`
	Longitude  float64   `json:"Longitude"`
	RecordedAt time.Time `json:"RecordedAt"`
}
