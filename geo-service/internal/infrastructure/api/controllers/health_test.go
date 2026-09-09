package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/geo-service/internal/core/ports/services"
	svcmocks "github.com/telemetry-platform/geo-service/internal/core/ports/services/mocks"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/api/controllers"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

func TestHealth_GetHealth(t *testing.T) {
	t.Parallel()

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()
		mockSvc := svcmocks.NewHealthService(t)
		mockSvc.On("GetHealth", mock.Anything).Return(services.HealthResult{Status: "ok", Service: "geo-service"}, nil)

		c := controllers.NewHealth(logger.NewLogger(), mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		c.GetHealth(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"status":"ok"`)
		assert.Contains(t, rec.Body.String(), `"service":"geo-service"`)
	})

	t.Run("fails when service fails", func(t *testing.T) {
		t.Parallel()
		mockSvc := svcmocks.NewHealthService(t)
		mockSvc.On("GetHealth", mock.Anything).Return(services.HealthResult{}, errors.New("internal error"))

		c := controllers.NewHealth(logger.NewLogger(), mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		c.GetHealth(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
