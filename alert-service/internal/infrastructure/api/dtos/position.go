package dtos

import "time"

// PositionSSEEvent represents a position event sent via Server-Sent Events.
type PositionSSEEvent struct {
	VehicleID  int64     `json:"vehicleId"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	RecordedAt time.Time `json:"recordedAt"`
}
