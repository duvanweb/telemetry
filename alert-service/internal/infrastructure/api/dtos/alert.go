package dtos

import "time"

// AlertSSEEvent represents an alert event sent via Server-Sent Events.
type AlertSSEEvent struct {
	ID         int64     `json:"id"`
	VehicleID  int64     `json:"vehicle_id"`
	Type       string    `json:"type"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	DetectedAt time.Time `json:"detected_at"`
}
