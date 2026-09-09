package rabbitmq

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	svcmocks "github.com/telemetry-platform/alert-service/internal/core/ports/services/mocks"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
	testdata "github.com/telemetry-platform/alert-service/test/data"
)

// mockDelivery is a test double for the Delivery interface.
type mockDelivery struct {
	acked   bool
	nacked  bool
	ackErr  error
	nackErr error
}

func (m *mockDelivery) Ack(bool) error       { m.acked = true; return m.ackErr }
func (m *mockDelivery) Nack(bool, bool) error { m.nacked = true; return m.nackErr }

func newTestConsumer(t *testing.T) (*Consumer, *svcmocks.AlertProcessor) {
	t.Helper()
	proc := svcmocks.NewAlertProcessor(t)
	return &Consumer{processor: proc, logger: logger.NewLogger()}, proc
}

func TestConsumer_ProcessMessage(t *testing.T) {
	t.Parallel()

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()
		c, proc := newTestConsumer(t)
		delivery := &mockDelivery{}

		pos := testdata.GetTestPosition()
		body, err := json.Marshal(pos)
		assert.NoError(t, err)

		proc.On("Process", mock.Anything, pos).Return(nil).Once()

		c.processMessage(context.Background(), body, delivery)

		assert.True(t, delivery.acked)
		assert.False(t, delivery.nacked)
	})

	t.Run("fails when process fails", func(t *testing.T) {
		t.Parallel()
		c, proc := newTestConsumer(t)
		delivery := &mockDelivery{}

		pos := testdata.GetTestPosition()
		body, err := json.Marshal(pos)
		assert.NoError(t, err)

		proc.On("Process", mock.Anything, pos).Return(errors.New("processing error")).Once()

		c.processMessage(context.Background(), body, delivery)

		assert.False(t, delivery.acked)
		assert.True(t, delivery.nacked)
	})

	t.Run("fails when unmarshal fails", func(t *testing.T) {
		t.Parallel()
		c, _ := newTestConsumer(t)
		delivery := &mockDelivery{}

		c.processMessage(context.Background(), []byte("invalid json"), delivery)

		assert.False(t, delivery.acked)
		assert.True(t, delivery.nacked)
	})
}

// Ensure domain import is used.
var _ = domain.AlertTypeVehicleStopped
