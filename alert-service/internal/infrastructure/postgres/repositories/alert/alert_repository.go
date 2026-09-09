package alert

import (
	"context"
	"fmt"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/core/ports/repositories"
)

// Compile-time check that Repository implements AlertRepository.
var _ repositories.AlertRepository = (*Repository)(nil)

// Repository implements repositories.AlertRepository with PostgreSQL.
type Repository struct {
	db repositories.Databaser
}

// NewRepository creates and returns a new alert Repository.
func NewRepository(db repositories.Databaser) *Repository {
	return &Repository{db: db}
}

// Save inserts a new alert into the database and populates the generated ID and CreatedAt.
func (r *Repository) Save(ctx context.Context, alert *domain.Alert) error {
	err := r.db.QueryRowContext(ctx, SaveAlertQuery,
		alert.VehicleID, alert.Type, alert.Latitude, alert.Longitude, alert.DetectedAt).
		Scan(&alert.ID, &alert.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save alert: %w", err)
	}
	return nil
}
