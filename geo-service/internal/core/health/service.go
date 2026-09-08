package health

import (
	"context"

	"github.com/telemetry-platform/geo-service/internal/core/ports/services"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// serviceName identifies this microservice in health responses.
const serviceName = "geo-service"

// Service implements the health check business logic.
type Service struct {
	logger logger.Logger
}

// GetHealth returns the current health status of the service.
func (s *Service) GetHealth(ctx context.Context) (services.HealthResult, error) {
	return services.HealthResult{Status: "ok", Service: serviceName}, nil
}

// NewService creates and returns a new health Service.
func NewService(log logger.Logger) *Service {
	return &Service{logger: log}
}
