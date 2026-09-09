package alert

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/alert-service/internal/core/ports/resources"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// Service implements the alert anomaly detection business logic.
type Service struct {
	repo            repositories.AlertRepository
	tracker         resources.PositionTracker
	broadcaster     resources.AlertBroadcaster
	stoppedThreshold time.Duration
	logger          logger.Logger
}

// NewService creates and returns a new alert Service.
// It parses the stopped threshold from the configuration string.
func NewService(
	repo repositories.AlertRepository,
	tracker resources.PositionTracker,
	broadcaster resources.AlertBroadcaster,
	config *env.Configuration,
	log logger.Logger,
) (*Service, error) {
	threshold, err := time.ParseDuration(config.StoppedThreshold)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stopped threshold: %w", err)
	}

	return &Service{
		repo:             repo,
		tracker:          tracker,
		broadcaster:      broadcaster,
		stoppedThreshold: threshold,
		logger:           log,
	}, nil
}

// Process evaluates a GPS position and detects the "Vehicle Stopped" anomaly.
// If the vehicle sends the same coordinates for longer than the stopped threshold,
// an alert is created, persisted, and broadcast to SSE subscribers.
func (s *Service) Process(ctx context.Context, pos domain.Position) error {
	track, err := s.tracker.Get(ctx, pos.VehicleID)
	if err != nil {
		if errors.Is(err, domain.ErrTrackNotFound) {
			return s.handleNewTrack(ctx, pos)
		}
		s.logger.Errorw(ctx, "failed to get vehicle track", "vehicle_id", pos.VehicleID, "error", err)
		return fmt.Errorf("failed to get vehicle track: %w", err)
	}

	if isSamePosition(track, pos) {
		return s.handleSamePosition(ctx, pos, track)
	}

	return s.handleMovedPosition(ctx, pos)
}

// handleNewTrack creates and stores a new vehicle track for the first seen position.
func (s *Service) handleNewTrack(ctx context.Context, pos domain.Position) error {
	newTrack := domain.VehicleTrack{
		VehicleID:   pos.VehicleID,
		Latitude:    pos.Latitude,
		Longitude:   pos.Longitude,
		FirstSeenAt: pos.RecordedAt,
		Alerted:     false,
	}

	if err := s.tracker.Set(ctx, newTrack); err != nil {
		s.logger.Errorw(ctx, "failed to set vehicle track", "vehicle_id", pos.VehicleID, "error", err)
		return fmt.Errorf("failed to set vehicle track: %w", err)
	}

	return nil
}

// handleMovedPosition resets the vehicle track when the vehicle moves to a new position.
func (s *Service) handleMovedPosition(ctx context.Context, pos domain.Position) error {
	newTrack := domain.VehicleTrack{
		VehicleID:   pos.VehicleID,
		Latitude:    pos.Latitude,
		Longitude:   pos.Longitude,
		FirstSeenAt: pos.RecordedAt,
		Alerted:     false,
	}

	if err := s.tracker.Set(ctx, newTrack); err != nil {
		s.logger.Errorw(ctx, "failed to set vehicle track", "vehicle_id", pos.VehicleID, "error", err)
		return fmt.Errorf("failed to set vehicle track: %w", err)
	}

	return nil
}

// handleSamePosition evaluates whether a stopped vehicle should trigger an alert.
func (s *Service) handleSamePosition(ctx context.Context, pos domain.Position, track domain.VehicleTrack) error {
	if track.Alerted {
		return nil
	}

	if pos.RecordedAt.Sub(track.FirstSeenAt) <= s.stoppedThreshold {
		return nil
	}

	alert := domain.Alert{
		VehicleID:  pos.VehicleID,
		Type:       domain.AlertTypeVehicleStopped,
		Latitude:   pos.Latitude,
		Longitude:  pos.Longitude,
		DetectedAt: pos.RecordedAt,
	}

	if err := s.repo.Save(ctx, alert); err != nil {
		s.logger.Errorw(ctx, "failed to save alert", "vehicle_id", pos.VehicleID, "error", err)
		return fmt.Errorf("failed to save alert: %w", err)
	}

	if err := s.broadcaster.Broadcast(ctx, alert); err != nil {
		s.logger.Errorw(ctx, "failed to broadcast alert", "vehicle_id", pos.VehicleID, "error", err)
	}

	track.Alerted = true
	if err := s.tracker.Set(ctx, track); err != nil {
		s.logger.Errorw(ctx, "failed to set vehicle track", "vehicle_id", pos.VehicleID, "error", err)
		return fmt.Errorf("failed to set vehicle track: %w", err)
	}

	return nil
}

// isSamePosition checks if the position matches the tracked coordinates.
func isSamePosition(track domain.VehicleTrack, pos domain.Position) bool {
	return track.Latitude == pos.Latitude && track.Longitude == pos.Longitude
}
