package services

import "context"

// HealthResult is the result of a health check.
type HealthResult struct {
	Status  string
	Service string
}

// HealthService is the interface for health checking.
//
//go:generate mockery --name HealthService --dir=. --output=./mocks
type HealthService interface {
	GetHealth(ctx context.Context) (HealthResult, error)
}
