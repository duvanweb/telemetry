package broadcaster_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/broadcaster"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
	testdata "github.com/telemetry-platform/alert-service/test/data"
)

func TestPositionBroadcaster_Broadcast(t *testing.T) {
	t.Parallel()

	t.Run("works correctly with subscriber", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewPositionBroadcaster(logger.NewLogger())
		ch := b.Subscribe()

		pos := testdata.GetTestPosition()
		err := b.Broadcast(context.Background(), pos)
		assert.NoError(t, err)

		received := <-ch
		assert.Equal(t, pos, received)

		b.Unsubscribe(ch)
	})

	t.Run("works correctly with no subscribers", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewPositionBroadcaster(logger.NewLogger())

		err := b.Broadcast(context.Background(), testdata.GetTestPosition())
		assert.NoError(t, err)
	})

	t.Run("handles correctly when subscriber buffer is full", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewPositionBroadcaster(logger.NewLogger())
		ch := b.Subscribe()

		// Fill the buffer (subscriberBufferSize = 16).
		for i := 0; i < 16; i++ {
			_ = b.Broadcast(context.Background(), testdata.GetTestPosition())
		}

		// Next broadcast should drop the position (non-blocking) but not error.
		err := b.Broadcast(context.Background(), testdata.GetTestPosition())
		assert.NoError(t, err)

		b.Unsubscribe(ch)
	})
}

func TestPositionBroadcaster_Subscribe(t *testing.T) {
	t.Parallel()

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewPositionBroadcaster(logger.NewLogger())
		ch := b.Subscribe()

		assert.NotNil(t, ch)

		b.Unsubscribe(ch)
	})
}

func TestPositionBroadcaster_Unsubscribe(t *testing.T) {
	t.Parallel()

	t.Run("works correctly and closes channel", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewPositionBroadcaster(logger.NewLogger())
		ch := b.Subscribe()

		b.Unsubscribe(ch)

		_, ok := <-ch
		assert.False(t, ok, "channel should be closed")
	})

	t.Run("handles correctly when channel not found", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewPositionBroadcaster(logger.NewLogger())
		other := make(chan domain.Position, 1)

		// Should not panic.
		b.Unsubscribe(other)
	})
}
