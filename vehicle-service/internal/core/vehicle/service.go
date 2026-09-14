package vehicle

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/telemetry-platform/vehicle-service/internal/core/domain"
	"github.com/telemetry-platform/vehicle-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/vehicle-service/internal/core/ports/resources"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
)

const (
	defaultLimit  = 20
	maxLimit      = 100
	defaultOffset = 0
)

// plateRegex validates the vehicle plate format: ABC-123.
var plateRegex = regexp.MustCompile(`^[A-Z]{3}-[0-9]{3}$`)

// Service implements the vehicle business logic.
type Service struct {
	repo          repositories.VehicleRepository
	eventPublisher resources.EventPublisher
	logger        logger.Logger
}

// Create validates the plate format and persists a new vehicle.
// Returns ErrInvalidPlate if the plate does not match ABC-123 format.
// Returns ErrVehicleAlreadyExists if an active vehicle with the same plate exists.
func (s *Service) Create(ctx context.Context, v *domain.Vehicle) error {
	if !plateRegex.MatchString(v.Plate) {
		return domain.ErrInvalidPlate
	}

	existing, err := s.repo.FindByPlate(ctx, v.Plate)
	if err != nil && !errors.Is(err, domain.ErrVehicleNotFound) {
		s.logger.Errorw(ctx, "failed to check existing vehicle by plate", "plate", v.Plate, "error", err)
		return err
	}
	if existing != nil {
		return domain.ErrVehicleAlreadyExists
	}

	if err := s.repo.Create(ctx, v); err != nil {
		s.logger.Errorw(ctx, "failed to create vehicle", "plate", v.Plate, "error", err)
		return err
	}

	return nil
}

// FindByPlate retrieves an active vehicle by its plate.
// Returns ErrVehicleNotFound if the vehicle does not exist or is soft-deleted.
func (s *Service) FindByPlate(ctx context.Context, plate string) (*domain.Vehicle, error) {
	v, err := s.repo.FindByPlate(ctx, plate)
	if err != nil {
		s.logger.Errorw(ctx, "failed to find vehicle by plate", "plate", plate, "error", err)
		return nil, err
	}

	return v, nil
}

// GetByID retrieves an active vehicle by its ID.
// Returns ErrVehicleNotFound if the vehicle does not exist or is soft-deleted.
func (s *Service) GetByID(ctx context.Context, id int64) (*domain.Vehicle, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw(ctx, "failed to get vehicle by id", "id", id, "error", err)
		return nil, err
	}

	return v, nil
}

// List retrieves a paginated list of active vehicles.
// Applies defaults limit=20/offset=0 and caps limit at 100.
func (s *Service) List(ctx context.Context, limit, offset int) ([]domain.Vehicle, int64, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = defaultOffset
	}

	vehicles, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Errorw(ctx, "failed to list vehicles", "limit", limit, "offset", offset, "error", err)
		return nil, 0, err
	}

	return vehicles, total, nil
}

// NewService creates and returns a new vehicle Service.
func NewService(repo repositories.VehicleRepository, pub resources.EventPublisher, log logger.Logger) *Service {
	return &Service{repo: repo, eventPublisher: pub, logger: log}
}

// SoftDelete marks a vehicle as deleted by setting deleted_at, then publishes a
// vehicle.deleted event so downstream services (geo-service, alert-service) can
// clean up orphaned data. If the event publish fails, the error is logged but
// NOT returned — the vehicle is already soft-deleted and the event will be
// retried via a reconciliation process or manual intervention.
func (s *Service) SoftDelete(ctx context.Context, id int64) error {
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		s.logger.Errorw(ctx, "failed to soft delete vehicle", "id", id, "error", err)
		return err
	}

	// Publish the vehicle.deleted event for downstream cleanup.
	// Best-effort: log on failure but don't fail the request.
	if err := s.eventPublisher.PublishVehicleDeleted(ctx, id, time.Now()); err != nil {
		s.logger.Errorw(ctx, "failed to publish vehicle deleted event", "id", id, "error", err)
	}

	return nil
}
