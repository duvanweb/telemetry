package vehicle

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/vehicle-service/internal/core/domain"
	"github.com/telemetry-platform/vehicle-service/internal/core/ports/repositories/mocks"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
)

// newTestService creates a vehicle Service with a mock repository for testing.
func newTestService(t *testing.T) (*Service, *mocks.VehicleRepository) {
	t.Helper()
	repo := mocks.NewVehicleRepository(t)
	svc := NewService(repo, logger.NewLogger())
	return svc, repo
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name        string
		plate       string
		setup       func(*mocks.VehicleRepository)
		expectedErr error
	}{
		{
			name:  "success",
			plate: "ABC-123",
			setup: func(m *mocks.VehicleRepository) {
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(nil, domain.ErrVehicleNotFound)
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:        "invalid plate lowercase",
			plate:       "abc-123",
			setup:       func(m *mocks.VehicleRepository) {},
			expectedErr: domain.ErrInvalidPlate,
		},
		{
			name:        "invalid plate too short",
			plate:       "ABC-12",
			setup:       func(m *mocks.VehicleRepository) {},
			expectedErr: domain.ErrInvalidPlate,
		},
		{
			name:        "invalid plate too long",
			plate:       "ABCD-123",
			setup:       func(m *mocks.VehicleRepository) {},
			expectedErr: domain.ErrInvalidPlate,
		},
		{
			name:  "already exists",
			plate: "ABC-123",
			setup: func(m *mocks.VehicleRepository) {
				existing := &domain.Vehicle{ID: 1, Plate: "ABC-123"}
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(existing, nil)
			},
			expectedErr: domain.ErrVehicleAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			tt.setup(repo)

			v := &domain.Vehicle{Plate: tt.plate}
			err := svc.Create(context.Background(), v)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestService_SoftDelete(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		setup       func(*mocks.VehicleRepository)
		expectedErr error
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *mocks.VehicleRepository) {
				m.On("SoftDelete", mock.Anything, int64(1)).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "not found",
			id:   999,
			setup: func(m *mocks.VehicleRepository) {
				m.On("SoftDelete", mock.Anything, int64(999)).Return(domain.ErrVehicleNotFound)
			},
			expectedErr: domain.ErrVehicleNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			tt.setup(repo)

			err := svc.SoftDelete(context.Background(), tt.id)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		id             int64
		setup          func(*mocks.VehicleRepository) *domain.Vehicle
		expectedErr    error
		expectVehicle  bool
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *mocks.VehicleRepository) *domain.Vehicle {
				v := &domain.Vehicle{ID: 1, Plate: "ABC-123", CreatedAt: time.Now(), UpdatedAt: time.Now()}
				m.On("GetByID", mock.Anything, int64(1)).Return(v, nil)
				return v
			},
			expectedErr:   nil,
			expectVehicle: true,
		},
		{
			name: "not found",
			id:   999,
			setup: func(m *mocks.VehicleRepository) *domain.Vehicle {
				m.On("GetByID", mock.Anything, int64(999)).Return(nil, domain.ErrVehicleNotFound)
				return nil
			},
			expectedErr:   domain.ErrVehicleNotFound,
			expectVehicle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
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

func TestService_FindByPlate(t *testing.T) {
	tests := []struct {
		name          string
		plate         string
		setup         func(*mocks.VehicleRepository) *domain.Vehicle
		expectedErr   error
		expectVehicle bool
	}{
		{
			name:  "success",
			plate: "ABC-123",
			setup: func(m *mocks.VehicleRepository) *domain.Vehicle {
				v := &domain.Vehicle{ID: 1, Plate: "ABC-123", CreatedAt: time.Now(), UpdatedAt: time.Now()}
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(v, nil)
				return v
			},
			expectedErr:   nil,
			expectVehicle: true,
		},
		{
			name:  "not found",
			plate: "ZZZ-999",
			setup: func(m *mocks.VehicleRepository) *domain.Vehicle {
				m.On("FindByPlate", mock.Anything, "ZZZ-999").Return(nil, domain.ErrVehicleNotFound)
				return nil
			},
			expectedErr:   domain.ErrVehicleNotFound,
			expectVehicle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
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

func TestService_List(t *testing.T) {
	now := time.Now()
	vehicles := []domain.Vehicle{
		{ID: 1, Plate: "ABC-123", CreatedAt: now, UpdatedAt: now},
		{ID: 2, Plate: "DEF-456", CreatedAt: now, UpdatedAt: now},
	}

	tests := []struct {
		name           string
		inputLimit     int
		inputOffset    int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "success with explicit params",
			inputLimit:     10,
			inputOffset:    5,
			expectedLimit:  10,
			expectedOffset: 5,
		},
		{
			name:           "defaults applied",
			inputLimit:     0,
			inputOffset:    0,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "max limit capped",
			inputLimit:     200,
			inputOffset:    0,
			expectedLimit:  100,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newTestService(t)
			repo.On("List", mock.Anything, tt.expectedLimit, tt.expectedOffset).
				Return(vehicles, int64(2), nil)

			result, total, err := svc.List(context.Background(), tt.inputLimit, tt.inputOffset)

			assert.NoError(t, err)
			assert.Equal(t, vehicles, result)
			assert.Equal(t, int64(2), total)
		})
	}
}
