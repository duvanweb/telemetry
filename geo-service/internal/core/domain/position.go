package domain

import "time"

// Position is a GPS reading for a vehicle at a point in time.
type Position struct {
	VehicleID  int64
	Latitude   float64
	Longitude  float64
	RecordedAt time.Time
}

// Validate checks that latitude and longitude are within valid ranges.
// Latitude must be in [-90, 90] and longitude in [-180, 180].
func (p Position) Validate() error {
	if p.Latitude < -90 || p.Latitude > 90 {
		return ErrInvalidPosition
	}
	if p.Longitude < -180 || p.Longitude > 180 {
		return ErrInvalidPosition
	}
	return nil
}
