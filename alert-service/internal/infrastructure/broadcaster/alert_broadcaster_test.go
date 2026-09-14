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

func TestBroadcaster_Broadcast(t *testing.T) {
	t.Parallel()

	t.Run("works correctly with subscriber", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewBroadcaster(logger.NewLogger())
		ch := b.Subscribe()

		alert := testdata.GetTestAlert()
		err := b.Broadcast(context.Background(), alert)
		assert.NoError(t, err)

		received := <-ch
		assert.Equal(t, alert, received)

		b.Unsubscribe(ch)
	})

	t.Run("works correctly with no subscribers", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewBroadcaster(logger.NewLogger())

		err := b.Broadcast(context.Background(), testdata.GetTestAlert())
		assert.NoError(t, err)
	})

	t.Run("handles correctly when subscriber buffer is full", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewBroadcaster(logger.NewLogger())
		ch := b.Subscribe()

		// Fill the buffer (subscriberBufferSize = 16).
		for i := 0; i < 16; i++ {
			_ = b.Broadcast(context.Background(), testdata.GetTestAlert())
		}

		// Next broadcast should drop the alert (non-blocking) but not error.
		err := b.Broadcast(context.Background(), testdata.GetTestAlert())
		assert.NoError(t, err)

		b.Unsubscribe(ch)
	})
}

func TestBroadcaster_Subscribe(t *testing.T) {
	t.Parallel()

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewBroadcaster(logger.NewLogger())
		ch := b.Subscribe()

		assert.NotNil(t, ch)

		b.Unsubscribe(ch)
	})
}

func TestBroadcaster_Unsubscribe(t *testing.T) {
	t.Parallel()

	t.Run("works correctly and closes channel", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewBroadcaster(logger.NewLogger())
		ch := b.Subscribe()

		b.Unsubscribe(ch)

		_, ok := <-ch
		assert.False(t, ok, "channel should be closed")
	})

	t.Run("handles correctly when channel not found", func(t *testing.T) {
		t.Parallel()

		b := broadcaster.NewBroadcaster(logger.NewLogger())
		other := make(chan domain.Alert, 1)

		// Should not panic.
		b.Unsubscribe(other)
	})
}
