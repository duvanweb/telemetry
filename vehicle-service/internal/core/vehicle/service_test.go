package vehicle_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/vehicle-service/internal/core/domain"
	reposmocks "github.com/telemetry-platform/vehicle-service/internal/core/ports/repositories/mocks"
	resmocks "github.com/telemetry-platform/vehicle-service/internal/core/ports/resources/mocks"
	"github.com/telemetry-platform/vehicle-service/internal/core/vehicle"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
	testdata "github.com/telemetry-platform/vehicle-service/test/data"
)

// newTestService creates a vehicle Service with mock repository and event publisher for testing.
func newTestService(t *testing.T) (*vehicle.Service, *reposmocks.VehicleRepository, *resmocks.EventPublisher) {
	t.Helper()
	repo := reposmocks.NewVehicleRepository(t)
	pub := resmocks.NewEventPublisher(t)
	svc := vehicle.NewService(repo, pub, logger.NewLogger())
	return svc, repo, pub
}

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		plate       string
		setup       func(*reposmocks.VehicleRepository)
		expectedErr error
	}{
		{
			name:  "works correctly",
			plate: "ABC-123",
			setup: func(m *reposmocks.VehicleRepository) {
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(nil, domain.ErrVehicleNotFound)
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:        "handles correctly when plate is invalid lowercase",
			plate:       "abc-123",
			setup:       func(m *reposmocks.VehicleRepository) {},
			expectedErr: domain.ErrInvalidPlate,
		},
		{
			name:        "handles correctly when plate is invalid too short",
			plate:       "ABC-12",
			setup:       func(m *reposmocks.VehicleRepository) {},
			expectedErr: domain.ErrInvalidPlate,
		},
		{
			name:        "handles correctly when plate is invalid too long",
			plate:       "ABCD-123",
			setup:       func(m *reposmocks.VehicleRepository) {},
			expectedErr: domain.ErrInvalidPlate,
		},
		{
			name:  "handles correctly when vehicle already exists",
			plate: "ABC-123",
			setup: func(m *reposmocks.VehicleRepository) {
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(testdata.GetTestVehicle(), nil)
			},
			expectedErr: domain.ErrVehicleAlreadyExists,
		},
		{
			name:  "fails when repository FindByPlate fails",
			plate: "ABC-123",
			setup: func(m *reposmocks.VehicleRepository) {
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(nil, assert.AnError)
			},
			expectedErr: assert.AnError,
		},
		{
			name:  "fails when repository Create fails",
			plate: "ABC-123",
			setup: func(m *reposmocks.VehicleRepository) {
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(nil, domain.ErrVehicleNotFound)
				m.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, repo, _ := newTestService(t)
			tt.setup(repo)

			v := &domain.Vehicle{Plate: tt.plate}
			err := svc.Create(context.Background(), v)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestService_FindByPlate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		plate         string
		setup         func(*reposmocks.VehicleRepository) *domain.Vehicle
		expectedErr   error
		expectVehicle bool
	}{
		{
			name:  "works correctly",
			plate: "ABC-123",
			setup: func(m *reposmocks.VehicleRepository) *domain.Vehicle {
				v := testdata.GetTestVehicle()
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(v, nil)
				return v
			},
			expectedErr:   nil,
			expectVehicle: true,
		},
		{
			name:  "handles correctly when vehicle not found",
			plate: "ZZZ-999",
			setup: func(m *reposmocks.VehicleRepository) *domain.Vehicle {
				m.On("FindByPlate", mock.Anything, "ZZZ-999").Return(nil, domain.ErrVehicleNotFound)
				return nil
			},
			expectedErr:   domain.ErrVehicleNotFound,
			expectVehicle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, repo, _ := newTestService(t)
			expected := tt.setup(repo)

			v, err := svc.FindByPlate(context.Background(), tt.plate)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectVehicle {
				assert.Equal(t, expected, v)
			} else {
				assert.Nil(t, v)
			}
		})
	}
}

func TestService_GetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		id            int64
		setup         func(*reposmocks.VehicleRepository) *domain.Vehicle
		expectedErr   error
		expectVehicle bool
	}{
		{
			name: "works correctly",
			id:   testdata.GetTestVehicleID(),
			setup: func(m *reposmocks.VehicleRepository) *domain.Vehicle {
				v := testdata.GetTestVehicle()
				m.On("GetByID", mock.Anything, testdata.GetTestVehicleID()).Return(v, nil)
				return v
			},
			expectedErr:   nil,
			expectVehicle: true,
		},
		{
			name: "handles correctly when vehicle not found",
			id:   999,
			setup: func(m *reposmocks.VehicleRepository) *domain.Vehicle {
				m.On("GetByID", mock.Anything, int64(999)).Return(nil, domain.ErrVehicleNotFound)
				return nil
			},
			expectedErr:   domain.ErrVehicleNotFound,
			expectVehicle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, repo, _ := newTestService(t)
			expected := tt.setup(repo)

			v, err := svc.GetByID(context.Background(), tt.id)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectVehicle {
				assert.Equal(t, expected, v)
			} else {
				assert.Nil(t, v)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	vehicles := testdata.GetTestVehicles()

	tests := []struct {
		name           string
		inputLimit     int
		inputOffset    int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "works correctly with explicit params",
			inputLimit:     10,
			inputOffset:    5,
			expectedLimit:  10,
			expectedOffset: 5,
		},
		{
			name:           "works correctly when defaults applied",
			inputLimit:     0,
			inputOffset:    0,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "handles correctly when max limit capped",
			inputLimit:     200,
			inputOffset:    0,
			expectedLimit:  100,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, repo, _ := newTestService(t)
			repo.On("List", mock.Anything, tt.expectedLimit, tt.expectedOffset).
				Return(vehicles, int64(len(vehicles)), nil)

			result, total, err := svc.List(context.Background(), tt.inputLimit, tt.inputOffset)

			assert.NoError(t, err)
			assert.Equal(t, vehicles, result)
			assert.Equal(t, int64(len(vehicles)), total)
		})
	}
}

func TestService_SoftDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		id          int64
		setup       func(*reposmocks.VehicleRepository, *resmocks.EventPublisher)
		expectedErr error
	}{
		{
			name: "works correctly",
			id:   testdata.GetTestVehicleID(),
			setup: func(m *reposmocks.VehicleRepository, p *resmocks.EventPublisher) {
				m.On("SoftDelete", mock.Anything, testdata.GetTestVehicleID()).Return(nil)
				p.On("PublishVehicleDeleted", mock.Anything, testdata.GetTestVehicleID(), mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "handles correctly when vehicle not found",
			id:   999,
			setup: func(m *reposmocks.VehicleRepository, _ *resmocks.EventPublisher) {
				m.On("SoftDelete", mock.Anything, int64(999)).Return(domain.ErrVehicleNotFound)
			},
			expectedErr: domain.ErrVehicleNotFound,
		},
		{
			name: "succeeds even when event publish fails",
			id:   testdata.GetTestVehicleID(),
			setup: func(m *reposmocks.VehicleRepository, p *resmocks.EventPublisher) {
				m.On("SoftDelete", mock.Anything, testdata.GetTestVehicleID()).Return(nil)
				p.On("PublishVehicleDeleted", mock.Anything, testdata.GetTestVehicleID(), mock.Anything).Return(assert.AnError)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, repo, pub := newTestService(t)
			tt.setup(repo, pub)

			err := svc.SoftDelete(context.Background(), tt.id)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
