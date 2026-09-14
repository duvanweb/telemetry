package dtos

// HealthResponse is the DTO for the health check response.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
