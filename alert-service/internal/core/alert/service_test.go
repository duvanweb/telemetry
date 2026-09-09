package alert_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/core/alert"
	repomocks "github.com/telemetry-platform/alert-service/internal/core/ports/repositories/mocks"
	resmocks "github.com/telemetry-platform/alert-service/internal/core/ports/resources/mocks"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// newTestService creates an alert Service with mock dependencies for testing.
func newTestService(t *testing.T) (*alert.Service, *repomocks.AlertRepository, *resmocks.PositionTracker, *resmocks.AlertBroadcaster) {
	t.Helper()
	repo := repomocks.NewAlertRepository(t)
	tracker := resmocks.NewPositionTracker(t)
	broadcaster := resmocks.NewAlertBroadcaster(t)

	config := &env.Configuration{StoppedThreshold: "60s"}
	svc, err := alert.NewService(repo, tracker, broadcaster, config, logger.NewLogger())
	require.NoError(t, err)

	return svc, repo, tracker, broadcaster
}

func TestService_Process(t *testing.T) {
	now := time.Now()
	lat := 4.71
	lng := -74.07

	t.Run("works correctly when first position", func(t *testing.T) {
		t.Parallel()
		svc, _, tracker, _ := newTestService(t)
		pos := domain.Position{VehicleID: 1, Latitude: lat, Longitude: lng, RecordedAt: now}

		tracker.On("Get", mock.Anything, int64(1)).Return(domain.VehicleTrack{}, domain.ErrTrackNotFound)
		tracker.On("Set", mock.Anything, mock.Anything).Return(nil)

		err := svc.Process(context.Background(), pos)
		assert.NoError(t, err)
	})

	t.Run("works correctly when same position below threshold", func(t *testing.T) {
		t.Parallel()
		svc, _, tracker, _ := newTestService(t)
		pos := domain.Position{VehicleID: 1, Latitude: lat, Longitude: lng, RecordedAt: now}
		track := domain.VehicleTrack{
			VehicleID:   1,
			Latitude:    lat,
			Longitude:   lng,
			FirstSeenAt: now.Add(-30 * time.Second),
			Alerted:     false,
		}

		tracker.On("Get", mock.Anything, int64(1)).Return(track, nil)

		err := svc.Process(context.Background(), pos)
		assert.NoError(t, err)
	})

	t.Run("works correctly when vehicle stopped detected", func(t *testing.T) {
		t.Parallel()
		svc, repo, tracker, broadcaster := newTestService(t)
		pos := domain.Position{VehicleID: 1, Latitude: lat, Longitude: lng, RecordedAt: now}
		track := domain.VehicleTrack{
			VehicleID:   1,
			Latitude:    lat,
			Longitude:   lng,
			FirstSeenAt: now.Add(-2 * time.Minute),
			Alerted:     false,
		}
		expectedAlert := domain.Alert{
			VehicleID:  1,
			Type:       domain.AlertTypeVehicleStopped,
			Latitude:   lat,
			Longitude:  lng,
			DetectedAt: now,
		}
		updatedTrack := track
		updatedTrack.Alerted = true

		tracker.On("Get", mock.Anything, int64(1)).Return(track, nil)
		repo.On("Save", mock.Anything, expectedAlert).Return(nil)
		broadcaster.On("Broadcast", mock.Anything, expectedAlert).Return(nil)
		tracker.On("Set", mock.Anything, updatedTrack).Return(nil)

		err := svc.Process(context.Background(), pos)
		assert.NoError(t, err)
	})

	t.Run("handles correctly when already alerted", func(t *testing.T) {
		t.Parallel()
		svc, _, tracker, _ := newTestService(t)
		pos := domain.Position{VehicleID: 1, Latitude: lat, Longitude: lng, RecordedAt: now}
		track := domain.VehicleTrack{
			VehicleID:   1,
			Latitude:    lat,
			Longitude:   lng,
			FirstSeenAt: now.Add(-2 * time.Minute),
			Alerted:     true,
		}

		tracker.On("Get", mock.Anything, int64(1)).Return(track, nil)

		err := svc.Process(context.Background(), pos)
		assert.NoError(t, err)
	})

	t.Run("handles correctly when vehicle moved", func(t *testing.T) {
		t.Parallel()
		svc, _, tracker, _ := newTestService(t)
		pos := domain.Position{VehicleID: 1, Latitude: 4.72, Longitude: -74.08, RecordedAt: now}
		track := domain.VehicleTrack{
			VehicleID:   1,
			Latitude:    lat,
			Longitude:   lng,
			FirstSeenAt: now.Add(-2 * time.Minute),
			Alerted:     false,
		}

		tracker.On("Get", mock.Anything, int64(1)).Return(track, nil)
		tracker.On("Set", mock.Anything, mock.Anything).Return(nil)

		err := svc.Process(context.Background(), pos)
		assert.NoError(t, err)
	})

	t.Run("fails when tracker get fails", func(t *testing.T) {
		t.Parallel()
		svc, _, tracker, _ := newTestService(t)
		pos := domain.Position{VehicleID: 1, Latitude: lat, Longitude: lng, RecordedAt: now}

		tracker.On("Get", mock.Anything, int64(1)).Return(domain.VehicleTrack{}, assert.AnError)

		err := svc.Process(context.Background(), pos)
		assert.Error(t, err)
	})

	t.Run("fails when tracker set fails on new track", func(t *testing.T) {
		t.Parallel()
		svc, _, tracker, _ := newTestService(t)
		pos := domain.Position{VehicleID: 1, Latitude: lat, Longitude: lng, RecordedAt: now}

		tracker.On("Get", mock.Anything, int64(1)).Return(domain.VehicleTrack{}, domain.ErrTrackNotFound)
		tracker.On("Set", mock.Anything, mock.Anything).Return(assert.AnError)

		err := svc.Process(context.Background(), pos)
		assert.Error(t, err)
	})

	t.Run("fails when repo save fails", func(t *testing.T) {
		t.Parallel()
		svc, repo, tracker, _ := newTestService(t)
		pos := domain.Position{VehicleID: 1, Latitude: lat, Longitude: lng, RecordedAt: now}
		track := domain.VehicleTrack{
			VehicleID:   1,
			Latitude:    lat,
			Longitude:   lng,
			FirstSeenAt: now.Add(-2 * time.Minute),
			Alerted:     false,
		}
		expectedAlert := domain.Alert{
			VehicleID:  1,
			Type:       domain.AlertTypeVehicleStopped,
			Latitude:   lat,
			Longitude:  lng,
			DetectedAt: now,
		}

		tracker.On("Get", mock.Anything, int64(1)).Return(track, nil)
		repo.On("Save", mock.Anything, expectedAlert).Return(assert.AnError)

		err := svc.Process(context.Background(), pos)
		assert.Error(t, err)
	})

	t.Run("fails when tracker set fails after alert", func(t *testing.T) {
		t.Parallel()
		svc, repo, tracker, broadcaster := newTestService(t)
		pos := domain.Position{VehicleID: 1, Latitude: lat, Longitude: lng, RecordedAt: now}
		track := domain.VehicleTrack{
			VehicleID:   1,
			Latitude:    lat,
			Longitude:   lng,
			FirstSeenAt: now.Add(-2 * time.Minute),
			Alerted:     false,
		}
		expectedAlert := domain.Alert{
			VehicleID:  1,
			Type:       domain.AlertTypeVehicleStopped,
			Latitude:   lat,
			Longitude:  lng,
			DetectedAt: now,
		}
		updatedTrack := track
		updatedTrack.Alerted = true

		tracker.On("Get", mock.Anything, int64(1)).Return(track, nil)
		repo.On("Save", mock.Anything, expectedAlert).Return(nil)
		broadcaster.On("Broadcast", mock.Anything, expectedAlert).Return(nil)
		tracker.On("Set", mock.Anything, updatedTrack).Return(assert.AnError)

		err := svc.Process(context.Background(), pos)
		assert.Error(t, err)
	})

	t.Run("fails when tracker set fails on moved position", func(t *testing.T) {
		t.Parallel()
		svc, _, tracker, _ := newTestService(t)
		pos := domain.Position{VehicleID: 1, Latitude: 4.72, Longitude: -74.08, RecordedAt: now}
		track := domain.VehicleTrack{
			VehicleID:   1,
			Latitude:    lat,
			Longitude:   lng,
			FirstSeenAt: now.Add(-2 * time.Minute),
			Alerted:     false,
		}

		tracker.On("Get", mock.Anything, int64(1)).Return(track, nil)
		tracker.On("Set", mock.Anything, mock.Anything).Return(assert.AnError)

		err := svc.Process(context.Background(), pos)
		assert.Error(t, err)
	})
}
