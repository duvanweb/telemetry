package vehicle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/telemetry-platform/vehicle-service/internal/core/domain"
	"github.com/telemetry-platform/vehicle-service/internal/core/ports/repositories"
)

// Compile-time check that Repository implements VehicleRepository.
var _ repositories.VehicleRepository = (*Repository)(nil)

// Repository implements repositories.VehicleRepository with PostgreSQL.
type Repository struct {
	db repositories.Databaser
}

// NewRepository creates and returns a new vehicle Repository.
func NewRepository(db repositories.Databaser) *Repository {
	return &Repository{db: db}
}

// Create inserts a new vehicle and populates the generated fields (ID, CreatedAt, UpdatedAt).
// Returns ErrVehicleAlreadyExists if a vehicle with the same plate already exists.
func (r *Repository) Create(ctx context.Context, v *domain.Vehicle) error {
	err := r.db.QueryRowContext(ctx, createVehicleQuery, v.Plate).
		Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrVehicleAlreadyExists
		}
		return fmt.Errorf("failed to create vehicle: %w", err)
	}
	return nil
}

// SoftDelete marks a vehicle as deleted by setting deleted_at and updated_at to NOW().
// Returns ErrVehicleNotFound if the vehicle does not exist or is already soft-deleted.
func (r *Repository) SoftDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, softDeleteVehicleQuery, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete vehicle: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrVehicleNotFound
	}
	return nil
}

// GetByID retrieves an active vehicle by its ID.
// Returns ErrVehicleNotFound if the vehicle does not exist or is soft-deleted.
func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.Vehicle, error) {
	v := &domain.Vehicle{}
	err := r.db.QueryRowContext(ctx, getVehicleByIDQuery, id).
		Scan(&v.ID, &v.Plate, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrVehicleNotFound
		}
		return nil, fmt.Errorf("failed to get vehicle by id: %w", err)
	}
	return v, nil
}

// FindByPlate retrieves an active vehicle by its plate.
// Returns ErrVehicleNotFound if the vehicle does not exist or is soft-deleted.
func (r *Repository) FindByPlate(ctx context.Context, plate string) (*domain.Vehicle, error) {
	v := &domain.Vehicle{}
	err := r.db.QueryRowContext(ctx, findVehicleByPlateQuery, plate).
		Scan(&v.ID, &v.Plate, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrVehicleNotFound
		}
		return nil, fmt.Errorf("failed to find vehicle by plate: %w", err)
	}
	return v, nil
}

// List retrieves a paginated list of active vehicles and the total count.
// Only vehicles with deleted_at IS NULL are returned.
func (r *Repository) List(ctx context.Context, limit, offset int) ([]domain.Vehicle, int64, error) {
	rows, err := r.db.QueryContext(ctx, listVehiclesQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list vehicles: %w", err)
	}
	defer func() { _ = rows.Close() }()

	vehicles := make([]domain.Vehicle, 0)
	for rows.Next() {
		var v domain.Vehicle
		if err := rows.Scan(&v.ID, &v.Plate, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, v)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, countVehiclesQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count vehicles: %w", err)
	}

	return vehicles, total, nil
}

// isUniqueViolation checks if the error is a PostgreSQL unique constraint violation.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == uniqueViolationCode
	}
	return false
}
