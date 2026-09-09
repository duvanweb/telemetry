package rabbitmq

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	svcmocks "github.com/telemetry-platform/geo-service/internal/core/ports/services/mocks"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// mockDelivery is a test double for the Delivery interface.
type mockDelivery struct {
	acked   bool
	nacked  bool
	ackErr  error
	nackErr error
}

func (m *mockDelivery) Ack(bool) error      { m.acked = true; return m.ackErr }
func (m *mockDelivery) Nack(bool, bool) error { m.nacked = true; return m.nackErr }

func newTestConsumer(t *testing.T) (*Consumer, *svcmocks.PositionProcessor) {
	t.Helper()
	proc := svcmocks.NewPositionProcessor(t)
	return &Consumer{processor: proc, logger: logger.NewLogger()}, proc
}

func TestConsumer_ProcessMessage_Success(t *testing.T) {
	c, proc := newTestConsumer(t)
	delivery := &mockDelivery{}

	pos := domain.Position{VehicleID: 1, Latitude: 4.71, Longitude: -74.07}
	body, err := json.Marshal(pos)
	assert.NoError(t, err)

	proc.On("Process", mock.Anything, pos).Return(nil).Once()

	c.processMessage(context.Background(), body, delivery)

	assert.True(t, delivery.acked)
	assert.False(t, delivery.nacked)
}

func TestConsumer_ProcessMessage_ProcessError(t *testing.T) {
	c, proc := newTestConsumer(t)
	delivery := &mockDelivery{}

	pos := domain.Position{VehicleID: 1, Latitude: 4.71, Longitude: -74.07}
	body, err := json.Marshal(pos)
	assert.NoError(t, err)

	proc.On("Process", mock.Anything, pos).Return(errors.New("db error")).Once()

	c.processMessage(context.Background(), body, delivery)

	assert.False(t, delivery.acked)
	assert.True(t, delivery.nacked)
}

func TestConsumer_ProcessMessage_CircuitBreakerOpen(t *testing.T) {
	c, proc := newTestConsumer(t)
	delivery := &mockDelivery{}

	pos := domain.Position{VehicleID: 1, Latitude: 4.71, Longitude: -74.07}
	body, err := json.Marshal(pos)
	assert.NoError(t, err)

	proc.On("Process", mock.Anything, pos).Return(domain.ErrCircuitBreakerOpen).Once()

	c.processMessage(context.Background(), body, delivery)

	assert.False(t, delivery.acked)
	assert.True(t, delivery.nacked)
}

func TestConsumer_ProcessMessage_UnmarshalError(t *testing.T) {
	c, _ := newTestConsumer(t)
	delivery := &mockDelivery{}

	c.processMessage(context.Background(), []byte("invalid json"), delivery)

	assert.False(t, delivery.acked)
	assert.True(t, delivery.nacked)
}
