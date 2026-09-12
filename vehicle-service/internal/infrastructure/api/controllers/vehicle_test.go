package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/vehicle-service/internal/core/domain"
	repomocks "github.com/telemetry-platform/vehicle-service/internal/core/ports/repositories/mocks"
	resmocks "github.com/telemetry-platform/vehicle-service/internal/core/ports/resources/mocks"
	"github.com/telemetry-platform/vehicle-service/internal/core/vehicle"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/api/controllers"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
	testdata "github.com/telemetry-platform/vehicle-service/test/data"
)

// newTestVehicleController creates a Vehicle controller with a mock repository for testing.
func newTestVehicleController(t *testing.T) (*controllers.Vehicle, *repomocks.VehicleRepository, *resmocks.EventPublisher) {
	t.Helper()
	repo := repomocks.NewVehicleRepository(t)
	pub := resmocks.NewEventPublisher(t)
	svc := vehicle.NewService(repo, pub, logger.NewLogger())
	c := controllers.NewVehicle(logger.NewLogger(), svc)
	return c, repo, pub
}

// newRequestWithChiParam creates a request with chi URL params set.
func newRequestWithChiParam(method, target, body string, params map[string]string) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	return req
}

func TestVehicle_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		body         string
		setup        func(*repomocks.VehicleRepository)
		expectedCode int
	}{
		{
			name:         "works correctly",
			body:         `{"plate":"ABC-123"}`,
			setup: func(m *repomocks.VehicleRepository) {
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(nil, domain.ErrVehicleNotFound)
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "handles correctly when body is invalid",
			body:         `{"plate":"abc-123"}`,
			setup:        func(m *repomocks.VehicleRepository) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "handles correctly when vehicle already exists",
			body:         `{"plate":"ABC-123"}`,
			setup: func(m *repomocks.VehicleRepository) {
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(testdata.GetTestVehicle(), nil)
			},
			expectedCode: http.StatusConflict,
		},
		{
			name:         "handles correctly when body is malformed",
			body:         `{invalid}`,
			setup:        func(m *repomocks.VehicleRepository) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, repo, _ := newTestVehicleController(t)
			tt.setup(repo)

			req := httptest.NewRequest(http.MethodPost, "/api/vehicles", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			c.Create(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestVehicle_FindByPlate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		plate        string
		setup        func(*repomocks.VehicleRepository)
		expectedCode int
	}{
		{
			name:  "works correctly",
			plate: "ABC-123",
			setup: func(m *repomocks.VehicleRepository) {
				m.On("FindByPlate", mock.Anything, "ABC-123").Return(testdata.GetTestVehicle(), nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:  "handles correctly when vehicle not found",
			plate: "ZZZ-999",
			setup: func(m *repomocks.VehicleRepository) {
				m.On("FindByPlate", mock.Anything, "ZZZ-999").Return(nil, domain.ErrVehicleNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, repo, _ := newTestVehicleController(t)
			tt.setup(repo)

			req := newRequestWithChiParam(http.MethodGet, "/api/vehicles/plates/"+tt.plate, "", map[string]string{"plate": tt.plate})
			rec := httptest.NewRecorder()

			c.FindByPlate(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestVehicle_GetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		id           string
		setup        func(*repomocks.VehicleRepository)
		expectedCode int
	}{
		{
			name: "works correctly",
			id:   "1",
			setup: func(m *repomocks.VehicleRepository) {
				m.On("GetByID", mock.Anything, int64(1)).Return(testdata.GetTestVehicle(), nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "handles correctly when id is invalid",
			id:           "abc",
			setup:        func(m *repomocks.VehicleRepository) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "handles correctly when vehicle not found",
			id:   "999",
			setup: func(m *repomocks.VehicleRepository) {
				m.On("GetByID", mock.Anything, int64(999)).Return(nil, domain.ErrVehicleNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, repo, _ := newTestVehicleController(t)
			tt.setup(repo)

			req := newRequestWithChiParam(http.MethodGet, "/api/vehicles/"+tt.id, "", map[string]string{"id": tt.id})
			rec := httptest.NewRecorder()

			c.GetByID(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestVehicle_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		query        string
		setup        func(*repomocks.VehicleRepository)
		expectedCode int
	}{
		{
			name:  "works correctly",
			query: "?limit=10&offset=0",
			setup: func(m *repomocks.VehicleRepository) {
				m.On("List", mock.Anything, 10, 0).Return(testdata.GetTestVehicles(), int64(2), nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:  "works correctly when defaults applied",
			query: "",
			setup: func(m *repomocks.VehicleRepository) {
				m.On("List", mock.Anything, 20, 0).Return(testdata.GetTestVehicles(), int64(2), nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, repo, _ := newTestVehicleController(t)
			tt.setup(repo)

			req := httptest.NewRequest(http.MethodGet, "/api/vehicles"+tt.query, nil)
			rec := httptest.NewRecorder()

			c.List(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestVehicle_SoftDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		id           string
		setup        func(*repomocks.VehicleRepository, *resmocks.EventPublisher)
		expectedCode int
	}{
		{
			name: "works correctly",
			id:   "1",
			setup: func(m *repomocks.VehicleRepository, p *resmocks.EventPublisher) {
				m.On("SoftDelete", mock.Anything, int64(1)).Return(nil)
				p.On("PublishVehicleDeleted", mock.Anything, int64(1), mock.Anything).Return(nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "handles correctly when id is invalid",
			id:           "abc",
			setup:        func(m *repomocks.VehicleRepository, _ *resmocks.EventPublisher) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "handles correctly when vehicle not found",
			id:   "999",
			setup: func(m *repomocks.VehicleRepository, _ *resmocks.EventPublisher) {
				m.On("SoftDelete", mock.Anything, int64(999)).Return(domain.ErrVehicleNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, repo, pub := newTestVehicleController(t)
			tt.setup(repo, pub)

			req := newRequestWithChiParam(http.MethodDelete, "/api/vehicles/"+tt.id, "", map[string]string{"id": tt.id})
			rec := httptest.NewRecorder()

			c.SoftDelete(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

// ensure encoding/json is used for marshaling in tests
var _ = json.Marshal

// ensure errors is used
var _ = errors.New
