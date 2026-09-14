package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/vehicle-service/internal/core/ports/services"
	svcmocks "github.com/telemetry-platform/vehicle-service/internal/core/ports/services/mocks"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/api/controllers"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/api/dtos"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
)

func TestHealth_GetHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setup        func(*svcmocks.HealthService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "works correctly",
			setup: func(m *svcmocks.HealthService) {
				m.On("GetHealth", mock.Anything).Return(services.HealthResult{Status: "ok", Service: "vehicle-service"}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: mustMarshalJSON(t, dtos.HealthResponse{Status: "ok", Service: "vehicle-service"}),
		},
		{
			name: "fails when service returns error",
			setup: func(m *svcmocks.HealthService) {
				m.On("GetHealth", mock.Anything).Return(services.HealthResult{}, errors.New("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockSvc := svcmocks.NewHealthService(t)
			tt.setup(mockSvc)
			c := controllers.NewHealth(logger.NewLogger(), mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			req = req.WithContext(context.Background())
			rec := httptest.NewRecorder()

			c.GetHealth(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}

// mustMarshalJSON marshals a value to JSON or fails the test.
func mustMarshalJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	assert.NoError(t, err)
	return string(b)
}
