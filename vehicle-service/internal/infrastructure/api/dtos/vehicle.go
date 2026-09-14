package dtos

import "time"

// CreateVehicleRequest is the request body for POST /api/vehicles.
type CreateVehicleRequest struct {
	Plate string `json:"plate"`
}

// VehicleResponse is the response body for single-vehicle endpoints.
type VehicleResponse struct {
	ID        int64     `json:"id"`
	Plate     string    `json:"plate"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ListVehiclesResponse is the paginated response for GET /api/vehicles.
type ListVehiclesResponse struct {
	Data   []VehicleResponse `json:"data"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
	Total  int64             `json:"total"`
}
