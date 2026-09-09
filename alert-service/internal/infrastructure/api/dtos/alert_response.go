package dtos

import "time"

// AlertResponse represents an alert in the REST API response.
type AlertResponse struct {
	ID         int64     `json:"id"`
	VehicleID  int64     `json:"vehicleId"`
	Type       string    `json:"type"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	DetectedAt time.Time `json:"detectedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

// ListAlertsResponse represents a paginated list of alerts in the REST API response.
type ListAlertsResponse struct {
	Data   []AlertResponse `json:"data"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Total  int64           `json:"total"`
}
