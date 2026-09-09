package repositories

import (
	"context"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
)

// AlertRepository is the interface for persisting alerts.
//
//go:generate mockery --name AlertRepository --dir=. --output=./mocks
type AlertRepository interface {
	Save(ctx context.Context, alert *domain.Alert) error
}
