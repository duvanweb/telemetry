package repositories

import (
	"context"

	"github.com/telemetry-platform/vehicle-service/internal/core/domain"
)

// VehicleRepository is the port for vehicle persistence.
//
//go:generate mockery --name VehicleRepository --dir=. --output=./mocks
type VehicleRepository interface {
	Create(ctx context.Context, v *domain.Vehicle) error
	SoftDelete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*domain.Vehicle, error)
	FindByPlate(ctx context.Context, plate string) (*domain.Vehicle, error)
	List(ctx context.Context, limit, offset int) ([]domain.Vehicle, int64, error)
}
