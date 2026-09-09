package alert

const (
	// SaveAlertQuery inserts a new alert and returns the generated id and created_at.
	SaveAlertQuery = `INSERT INTO alerts (vehicle_id, type, latitude, longitude, detected_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	// ListAlertsQuery selects a page of alerts ordered by detected_at descending.
	ListAlertsQuery = `SELECT id, vehicle_id, type, latitude, longitude, detected_at, created_at FROM alerts ORDER BY detected_at DESC LIMIT $1 OFFSET $2`

	// CountAlertsQuery counts the total number of alerts.
	CountAlertsQuery = `SELECT COUNT(*) FROM alerts`
)
