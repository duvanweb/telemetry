package processor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	repomocks "github.com/telemetry-platform/geo-service/internal/core/ports/repositories/mocks"
	"github.com/telemetry-platform/geo-service/internal/core/processor"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

func newTestProcessor(t *testing.T, threshold uint32, timeout time.Duration) (*processor.Processor, *repomocks.PositionRepository) {
	t.Helper()
	repo := repomocks.NewPositionRepository(t)
	proc := processor.NewProcessor(repo, threshold, timeout, logger.NewLogger())
	return proc, repo
}

func TestProcessor_Process_Success(t *testing.T) {
	proc, repo := newTestProcessor(t, 5, 30*time.Second)

	pos := domain.Position{VehicleID: 1, Latitude: 4.71, Longitude: -74.07}
	repo.On("Save", mock.Anything, pos).Return(nil).Once()

	err := proc.Process(context.Background(), pos)
	require.NoError(t, err)
}

func TestProcessor_Process_CircuitBreakerOpen(t *testing.T) {
	proc, repo := newTestProcessor(t, 1, 30*time.Second)

	pos := domain.Position{VehicleID: 1, Latitude: 4.71, Longitude: -74.07}

	// First call fails — CB opens (threshold=1).
	repo.On("Save", mock.Anything, pos).Return(errors.New("db connection refused")).Once()

	err := proc.Process(context.Background(), pos)
	require.Error(t, err)
	assert.NotEqual(t, domain.ErrCircuitBreakerOpen, err)

	// Second call — CB is open, repo.Save is NOT called.
	err = proc.Process(context.Background(), pos)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrCircuitBreakerOpen)
}

func TestProcessor_Process_OpensAfterThreshold(t *testing.T) {
	proc, repo := newTestProcessor(t, 5, 30*time.Second)

	pos := domain.Position{VehicleID: 1, Latitude: 4.71, Longitude: -74.07}
	dbErr := errors.New("db connection refused")

	// 5 consecutive failures — CB should open after the 5th.
	repo.On("Save", mock.Anything, pos).Return(dbErr).Times(5)

	for i := 0; i < 5; i++ {
		err := proc.Process(context.Background(), pos)
		require.Error(t, err)
		assert.NotErrorIs(t, err, domain.ErrCircuitBreakerOpen, "call %d should not be CB open", i+1)
	}

	// 6th call — CB is open, repo.Save is NOT called.
	err := proc.Process(context.Background(), pos)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrCircuitBreakerOpen)
}
