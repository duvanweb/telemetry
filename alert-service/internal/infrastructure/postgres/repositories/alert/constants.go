package alert

const (
	// SaveAlertQuery inserts a new alert into the alerts table.
	SaveAlertQuery = `INSERT INTO alerts (vehicle_id, type, latitude, longitude, detected_at) VALUES ($1, $2, $3, $4, $5)`
)
