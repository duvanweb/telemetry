package tracker_test

import (
	"context"
	"testing"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/redis/tracker"
	testdata "github.com/telemetry-platform/alert-service/test/data"
)

// newTestTracker creates a tracker connected to a local Redis instance.
// Skips the test if Redis is not available.
func newTestTracker(t *testing.T) *tracker.Tracker {
	t.Helper()
	client := goredis.NewClient(&goredis.Options{Addr: "localhost:6379"})

	ctx, cancel := context.WithTimeout(context.Background(), 2)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping tracker test")
	}

	t.Cleanup(func() {
		_ = client.FlushAll(context.Background())
		_ = client.Close()
	})

	return tracker.NewTracker(client)
}

func TestTracker_Get(t *testing.T) {
	t.Parallel()

	trk := newTestTracker(t)

	t.Run("works correctly when track exists", func(t *testing.T) {
		t.Parallel()

		track := testdata.GetTestVehicleTrack()
		err := trk.Set(context.Background(), track)
		assert.NoError(t, err)

		result, err := trk.Get(context.Background(), track.VehicleID)
		assert.NoError(t, err)
		assert.Equal(t, track, result)
	})

	t.Run("handles correctly when track not found", func(t *testing.T) {
		t.Parallel()

		result, err := trk.Get(context.Background(), 99999)
		assert.Equal(t, domain.ErrTrackNotFound, err)
		assert.Equal(t, domain.VehicleTrack{}, result)
	})
}

func TestTracker_Set(t *testing.T) {
	t.Parallel()

	trk := newTestTracker(t)

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()

		track := testdata.GetTestVehicleTrack()
		err := trk.Set(context.Background(), track)
		assert.NoError(t, err)

		result, err := trk.Get(context.Background(), track.VehicleID)
		assert.NoError(t, err)
		assert.Equal(t, track, result)
	})

	t.Run("works correctly when overwriting existing track", func(t *testing.T) {
		t.Parallel()

		track := testdata.GetTestVehicleTrack()
		err := trk.Set(context.Background(), track)
		assert.NoError(t, err)

		updated := testdata.GetTestAlertedTrack()
		err = trk.Set(context.Background(), updated)
		assert.NoError(t, err)

		result, err := trk.Get(context.Background(), track.VehicleID)
		assert.NoError(t, err)
		assert.True(t, result.Alerted)
	})
}
