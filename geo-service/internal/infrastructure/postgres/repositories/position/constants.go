package position

const (
	// SavePositionQuery inserts a new GPS position into the gps_positions table.
	SavePositionQuery = `INSERT INTO gps_positions (vehicle_id, latitude, longitude, recorded_at) VALUES ($1, $2, $3, $4)`
)
