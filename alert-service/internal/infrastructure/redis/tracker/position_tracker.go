package tracker

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	goredis "github.com/redis/go-redis/v9"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/core/ports/resources"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Compile-time check that Tracker implements PositionTracker.
var _ resources.PositionTracker = (*Tracker)(nil)

// Tracker implements resources.PositionTracker with Redis.
type Tracker struct {
	client *goredis.Client
}

// NewTracker creates and returns a new position Tracker.
func NewTracker(client *goredis.Client) *Tracker {
	return &Tracker{client: client}
}

// Get retrieves the vehicle track from Redis.
// Returns domain.ErrTrackNotFound if the track does not exist.
func (t *Tracker) Get(ctx context.Context, vehicleID int64) (domain.VehicleTrack, error) {
	key := buildTrackKey(vehicleID)

	data, err := t.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return domain.VehicleTrack{}, domain.ErrTrackNotFound
		}
		return domain.VehicleTrack{}, fmt.Errorf("failed to get vehicle track: %w", err)
	}

	var track domain.VehicleTrack
	if err := json.Unmarshal(data, &track); err != nil {
		return domain.VehicleTrack{}, fmt.Errorf("failed to unmarshal vehicle track: %w", err)
	}

	return track, nil
}

// Set stores the vehicle track in Redis.
func (t *Tracker) Set(ctx context.Context, track domain.VehicleTrack) error {
	key := buildTrackKey(track.VehicleID)

	data, err := json.Marshal(track)
	if err != nil {
		return fmt.Errorf("failed to marshal vehicle track: %w", err)
	}

	if err := t.client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to set vehicle track: %w", err)
	}

	return nil
}

// Delete removes the vehicle track from Redis.
// This is invoked when a vehicle is deleted to clean up orphaned tracker data.
// The operation is idempotent — deleting a non-existent key is a no-op.
func (t *Tracker) Delete(ctx context.Context, vehicleID int64) error {
	key := buildTrackKey(vehicleID)
	if err := t.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete vehicle track: %w", err)
	}
	return nil
}

// buildTrackKey builds the Redis key for a vehicle track.
func buildTrackKey(vehicleID int64) string {
	return "alert:vehicle:" + strconv.FormatInt(vehicleID, 10)
}
