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

// List returns a page of alerts ordered by detected_at descending and the total count.
func (r *Repository) List(ctx context.Context, limit, offset int) ([]domain.Alert, int64, error) {
	rows, err := r.db.QueryContext(ctx, ListAlertsQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list alerts: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var alerts []domain.Alert
	for rows.Next() {
		var alert domain.Alert
		if err := rows.Scan(&alert.ID, &alert.VehicleID, &alert.Type, &alert.Latitude, &alert.Longitude, &alert.DetectedAt, &alert.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan alert: %w", err)
		}
		alerts = append(alerts, alert)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate alert rows: %w", err)
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, CountAlertsQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count alerts: %w", err)
	}

	return alerts, total, nil
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

// DeleteByVehicleID deletes all alerts for the given vehicle ID.
// This is invoked when a vehicle is deleted to clean up orphaned alert data.
// The operation is idempotent — deleting from an empty set is a no-op.
func (r *Repository) DeleteByVehicleID(ctx context.Context, vehicleID int64) error {
	_, err := r.db.ExecContext(ctx, DeleteByVehicleIDQuery, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to delete alerts by vehicle id: %w", err)
	}
	return nil
}
