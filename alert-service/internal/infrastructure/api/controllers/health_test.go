package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/alert-service/internal/core/ports/services"
	svcmocks "github.com/telemetry-platform/alert-service/internal/core/ports/services/mocks"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/controllers"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/dtos"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

func TestHealth_GetHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func(*svcmocks.HealthService)
		expectedCode   int
		expectedStatus string
	}{
		{
			name: "works correctly",
			setup: func(m *svcmocks.HealthService) {
				m.On("GetHealth", mock.Anything).Return(services.HealthResult{Status: "ok", Service: "alert-service"}, nil)
			},
			expectedCode:   http.StatusOK,
			expectedStatus: "ok",
		},
		{
			name: "fails when service fails",
			setup: func(m *svcmocks.HealthService) {
				m.On("GetHealth", mock.Anything).Return(services.HealthResult{}, assert.AnError)
			},
			expectedCode:   http.StatusInternalServerError,
			expectedStatus: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := svcmocks.NewHealthService(t)
			tt.setup(svc)

			c := controllers.NewHealth(logger.NewLogger(), svc)
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			c.GetHealth(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == http.StatusOK {
				var resp dtos.HealthResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.Status)
				assert.Equal(t, "alert-service", resp.Service)
			}
		})
	}
}
