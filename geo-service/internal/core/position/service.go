package position

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	"github.com/telemetry-platform/geo-service/internal/core/ports/resources"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// Service implements the position ingestion business logic.
type Service struct {
	vehicleClient resources.VehicleClient
	cache         resources.PositionCache
	publisher     resources.PositionPublisher
	cacheTTL      time.Duration
	logger        logger.Logger
}

// Ingest validates a position, checks for duplicates, caches it, and publishes it for async persistence.
// Returns ErrInvalidPosition if lat/lng are out of range.
// Returns ErrVehicleNotFound if the vehicle does not exist.
// Returns ErrVehicleServiceUnavailable if vehicle-service is unreachable.
// Returns ErrDuplicatePosition if the same position was recently ingested.
func (s *Service) Ingest(ctx context.Context, pos domain.Position) error {
	if err := pos.Validate(); err != nil {
		return err
	}

	if err := s.vehicleClient.ValidateVehicle(ctx, pos.VehicleID); err != nil {
		return err
	}

	key := buildCacheKey(pos)
	exists, err := s.cache.Exists(ctx, key)
	if err != nil {
		s.logger.Errorw(ctx, "failed to check cache for duplicate", "error", err)
		return fmt.Errorf("failed to check cache: %w", err)
	}
	if exists {
		return domain.ErrDuplicatePosition
	}

	if err := s.cache.Set(ctx, key, pos.RecordedAt.Format(time.RFC3339), s.cacheTTL); err != nil {
		s.logger.Errorw(ctx, "failed to set cache", "error", err)
		return fmt.Errorf("failed to set cache: %w", err)
	}

	if err := s.publisher.Publish(ctx, pos); err != nil {
		s.logger.Errorw(ctx, "failed to publish position", "error", err)
		return fmt.Errorf("failed to publish position: %w", err)
	}

	return nil
}

// NewService creates and returns a new position Service.
// It parses the cache TTL from the configuration string.
func NewService(
	vehicleClient resources.VehicleClient,
	cache resources.PositionCache,
	publisher resources.PositionPublisher,
	config *env.Configuration,
	log logger.Logger,
) (*Service, error) {
	ttl, err := time.ParseDuration(config.CacheTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cache TTL: %w", err)
	}

	return &Service{
		vehicleClient: vehicleClient,
		cache:         cache,
		publisher:     publisher,
		cacheTTL:      ttl,
		logger:        log,
	}, nil
}

// buildCacheKey builds the Redis cache key for a position.
// Format: geo:pos:{vehicleID}:{lat}:{lng}
func buildCacheKey(pos domain.Position) string {
	return fmt.Sprintf("geo:pos:%d:%s:%s",
		pos.VehicleID,
		strconv.FormatFloat(pos.Latitude, 'f', -1, 64),
		strconv.FormatFloat(pos.Longitude, 'f', -1, 64),
	)
}
