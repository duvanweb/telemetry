package alert

const (
	// SaveAlertQuery inserts a new alert and returns the generated id and created_at.
	SaveAlertQuery = `INSERT INTO alerts (vehicle_id, type, latitude, longitude, detected_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
)
