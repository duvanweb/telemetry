package alert

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/core/ports/repositories"
)

// Compile-time check that Repository implements AlertRepository.
var _ repositories.AlertRepository = (*Repository)(nil)

// Repository implements repositories.AlertRepository with PostgreSQL.
type Repository struct {
	db *sql.DB
}

// NewRepository creates and returns a new alert Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Save inserts a new alert into the database.
func (r *Repository) Save(ctx context.Context, alert domain.Alert) error {
	_, err := r.db.ExecContext(ctx, SaveAlertQuery,
		alert.VehicleID, alert.Type, alert.Latitude, alert.Longitude, alert.DetectedAt)
	if err != nil {
		return fmt.Errorf("failed to save alert: %w", err)
	}
	return nil
}
