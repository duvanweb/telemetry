package position

import (
	"context"
	"fmt"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	"github.com/telemetry-platform/geo-service/internal/core/ports/repositories"
)

// Compile-time check that Repository implements PositionRepository.
var _ repositories.PositionRepository = (*Repository)(nil)

// Repository implements repositories.PositionRepository with PostgreSQL.
type Repository struct {
	db repositories.Databaser
}

// NewRepository creates and returns a new position Repository.
func NewRepository(db repositories.Databaser) *Repository {
	return &Repository{db: db}
}

// Save inserts a new GPS position into the database.
func (r *Repository) Save(ctx context.Context, pos domain.Position) error {
	_, err := r.db.ExecContext(ctx, SavePositionQuery,
		pos.VehicleID, pos.Latitude, pos.Longitude, pos.RecordedAt)
	if err != nil {
		return fmt.Errorf("failed to save position: %w", err)
	}
	return nil
}
