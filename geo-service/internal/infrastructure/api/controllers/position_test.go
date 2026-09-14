package controllers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	svcmocks "github.com/telemetry-platform/geo-service/internal/core/ports/services/mocks"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/api/controllers"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
	testdata "github.com/telemetry-platform/geo-service/test/data"
)

func TestPosition_Create(t *testing.T) {
	t.Parallel()

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()
		mockSvc := svcmocks.NewPositionService(t)
		pos := testdata.GetTestPosition()
		mockSvc.On("Ingest", mock.Anything, pos).Return(nil)

		c := controllers.NewPosition(logger.NewLogger(), mockSvc)
		body := `{"lat":4.71,"lng":-74.07,"timestamp":"2026-01-01T12:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/api/vehicles/1/positions", strings.NewReader(body))
		req = setChiURLParam(req, "vehicle_id", "1")
		rec := httptest.NewRecorder()

		c.Create(rec, req)

		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Contains(t, rec.Body.String(), `"status":"accepted"`)
	})

	t.Run("fails when vehicle id is invalid", func(t *testing.T) {
		t.Parallel()
		mockSvc := svcmocks.NewPositionService(t)

		c := controllers.NewPosition(logger.NewLogger(), mockSvc)
		body := `{"lat":4.71,"lng":-74.07,"timestamp":"2026-01-01T12:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/api/vehicles/abc/positions", strings.NewReader(body))
		req = setChiURLParam(req, "vehicle_id", "abc")
		rec := httptest.NewRecorder()

		c.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("fails when body is invalid", func(t *testing.T) {
		t.Parallel()
		mockSvc := svcmocks.NewPositionService(t)

		c := controllers.NewPosition(logger.NewLogger(), mockSvc)
		req := httptest.NewRequest(http.MethodPost, "/api/vehicles/1/positions", strings.NewReader("invalid json"))
		req = setChiURLParam(req, "vehicle_id", "1")
		rec := httptest.NewRecorder()

		c.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("handles correctly when vehicle not found", func(t *testing.T) {
		t.Parallel()
		mockSvc := svcmocks.NewPositionService(t)
		pos := testdata.GetTestPosition()
		mockSvc.On("Ingest", mock.Anything, pos).Return(domain.ErrVehicleNotFound)

		c := controllers.NewPosition(logger.NewLogger(), mockSvc)
		body := `{"lat":4.71,"lng":-74.07,"timestamp":"2026-01-01T12:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/api/vehicles/1/positions", strings.NewReader(body))
		req = setChiURLParam(req, "vehicle_id", "1")
		rec := httptest.NewRecorder()

		c.Create(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("handles correctly when position is duplicate", func(t *testing.T) {
		t.Parallel()
		mockSvc := svcmocks.NewPositionService(t)
		pos := testdata.GetTestPosition()
		mockSvc.On("Ingest", mock.Anything, pos).Return(domain.ErrDuplicatePosition)

		c := controllers.NewPosition(logger.NewLogger(), mockSvc)
		body := `{"lat":4.71,"lng":-74.07,"timestamp":"2026-01-01T12:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/api/vehicles/1/positions", strings.NewReader(body))
		req = setChiURLParam(req, "vehicle_id", "1")
		rec := httptest.NewRecorder()

		c.Create(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("handles correctly when vehicle service unavailable", func(t *testing.T) {
		t.Parallel()
		mockSvc := svcmocks.NewPositionService(t)
		pos := testdata.GetTestPosition()
		mockSvc.On("Ingest", mock.Anything, pos).Return(domain.ErrVehicleServiceUnavailable)

		c := controllers.NewPosition(logger.NewLogger(), mockSvc)
		body := `{"lat":4.71,"lng":-74.07,"timestamp":"2026-01-01T12:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/api/vehicles/1/positions", strings.NewReader(body))
		req = setChiURLParam(req, "vehicle_id", "1")
		rec := httptest.NewRecorder()

		c.Create(rec, req)

		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	})
}

// setChiURLParam sets a chi URL parameter on the request for testing.
func setChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
