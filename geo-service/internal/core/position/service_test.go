package position_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	resmocks "github.com/telemetry-platform/geo-service/internal/core/ports/resources/mocks"
	"github.com/telemetry-platform/geo-service/internal/core/position"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
	testdata "github.com/telemetry-platform/geo-service/test/data"
)

// newTestService creates a position Service with mock dependencies for testing.
func newTestService(t *testing.T) (*position.Service, *resmocks.VehicleClient, *resmocks.PositionCache, *resmocks.PositionPublisher) {
	t.Helper()
	vehicleClient := resmocks.NewVehicleClient(t)
	cache := resmocks.NewPositionCache(t)
	publisher := resmocks.NewPositionPublisher(t)

	config := &env.Configuration{CacheTTL: "60s"}
	svc, err := position.NewService(vehicleClient, cache, publisher, config, logger.NewLogger())
	require.NoError(t, err)

	return svc, vehicleClient, cache, publisher
}

func TestService_Ingest(t *testing.T) {
	t.Parallel()

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()
		svc, vehicleClient, cache, publisher := newTestService(t)
		pos := testdata.GetTestPosition()

		vehicleClient.On("ValidateVehicle", mock.Anything, int64(1)).Return(nil)
		cache.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
		cache.On("Set", mock.Anything, mock.Anything, mock.Anything, 60*time.Second).Return(nil)
		publisher.On("Publish", mock.Anything, pos).Return(nil)

		err := svc.Ingest(context.Background(), pos)
		assert.NoError(t, err)
	})

	t.Run("handles correctly when position is invalid", func(t *testing.T) {
		t.Parallel()
		svc, _, _, _ := newTestService(t)
		pos := testdata.GetTestInvalidPosition()

		err := svc.Ingest(context.Background(), pos)
		assert.Equal(t, domain.ErrInvalidPosition, err)
	})

	t.Run("handles correctly when vehicle not found", func(t *testing.T) {
		t.Parallel()
		svc, vehicleClient, _, _ := newTestService(t)
		pos := testdata.GetTestPosition()

		vehicleClient.On("ValidateVehicle", mock.Anything, int64(1)).Return(domain.ErrVehicleNotFound)

		err := svc.Ingest(context.Background(), pos)
		assert.Equal(t, domain.ErrVehicleNotFound, err)
	})

	t.Run("handles correctly when vehicle service unavailable", func(t *testing.T) {
		t.Parallel()
		svc, vehicleClient, _, _ := newTestService(t)
		pos := testdata.GetTestPosition()

		vehicleClient.On("ValidateVehicle", mock.Anything, int64(1)).Return(domain.ErrVehicleServiceUnavailable)

		err := svc.Ingest(context.Background(), pos)
		assert.Equal(t, domain.ErrVehicleServiceUnavailable, err)
	})

	t.Run("handles correctly when position is duplicate", func(t *testing.T) {
		t.Parallel()
		svc, vehicleClient, cache, _ := newTestService(t)
		pos := testdata.GetTestPosition()

		vehicleClient.On("ValidateVehicle", mock.Anything, int64(1)).Return(nil)
		cache.On("Exists", mock.Anything, mock.Anything).Return(true, nil)

		err := svc.Ingest(context.Background(), pos)
		assert.Equal(t, domain.ErrDuplicatePosition, err)
	})

	t.Run("fails when cache exists check fails", func(t *testing.T) {
		t.Parallel()
		svc, vehicleClient, cache, _ := newTestService(t)
		pos := testdata.GetTestPosition()

		vehicleClient.On("ValidateVehicle", mock.Anything, int64(1)).Return(nil)
		cache.On("Exists", mock.Anything, mock.Anything).Return(false, assert.AnError)

		err := svc.Ingest(context.Background(), pos)
		assert.Error(t, err)
	})

	t.Run("fails when cache set fails", func(t *testing.T) {
		t.Parallel()
		svc, vehicleClient, cache, _ := newTestService(t)
		pos := testdata.GetTestPosition()

		vehicleClient.On("ValidateVehicle", mock.Anything, int64(1)).Return(nil)
		cache.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
		cache.On("Set", mock.Anything, mock.Anything, mock.Anything, 60*time.Second).Return(assert.AnError)

		err := svc.Ingest(context.Background(), pos)
		assert.Error(t, err)
	})

	t.Run("fails when publisher fails", func(t *testing.T) {
		t.Parallel()
		svc, vehicleClient, cache, publisher := newTestService(t)
		pos := testdata.GetTestPosition()

		vehicleClient.On("ValidateVehicle", mock.Anything, int64(1)).Return(nil)
		cache.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
		cache.On("Set", mock.Anything, mock.Anything, mock.Anything, 60*time.Second).Return(nil)
		publisher.On("Publish", mock.Anything, pos).Return(assert.AnError)

		err := svc.Ingest(context.Background(), pos)
		assert.Error(t, err)
	})
}
