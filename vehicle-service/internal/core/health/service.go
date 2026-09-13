package health

import (
	"context"
	"time"

	"github.com/telemetry-platform/vehicle-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/vehicle-service/internal/core/ports/services"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
)

// serviceName identifies this microservice in health responses.
const serviceName = "vehicle-service"

// Service implements the health check business logic.
type Service struct {
	db     repositories.Databaser
	logger logger.Logger
}

// GetHealth returns the current health status of the service.
// It pings the database with a 2s timeout. If the DB is unreachable,
// the status is "degraded" but the HTTP response is still 200 so Docker
// does not restart the service — the process is alive, only the DB is down.
func (s *Service) GetHealth(ctx context.Context) (services.HealthResult, error) {
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := s.db.PingContext(pingCtx); err != nil {
		s.logger.Warnw(ctx, "health check: database unreachable", "error", err)
		return services.HealthResult{Status: "degraded", Service: serviceName}, nil
	}

	return services.HealthResult{Status: "ok", Service: serviceName}, nil
}

// NewService creates and returns a new health Service.
func NewService(db repositories.Databaser, log logger.Logger) *Service {
	return &Service{db: db, logger: log}
}
