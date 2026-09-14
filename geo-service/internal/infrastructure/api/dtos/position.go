package dtos

// CreatePositionRequest is the JSON body for POST /api/vehicles/{vehicle_id}/positions.
type CreatePositionRequest struct {
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Timestamp string  `json:"timestamp"` // RFC3339
}

// CreatePositionResponse is the JSON body for 202 Accepted.
type CreatePositionResponse struct {
	Status    string `json:"status"`     // "accepted"
	VehicleID int64  `json:"vehicle_id"`
}
