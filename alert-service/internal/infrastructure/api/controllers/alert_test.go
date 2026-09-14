package controllers_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	repomocks "github.com/telemetry-platform/alert-service/internal/core/ports/repositories/mocks"
	resmocks "github.com/telemetry-platform/alert-service/internal/core/ports/resources/mocks"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/controllers"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/dtos"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
	testdata "github.com/telemetry-platform/alert-service/test/data"
)

// newTestAlertController creates an Alert controller with mock dependencies for testing.
func newTestAlertController(t *testing.T) (*controllers.Alert, *repomocks.AlertRepository, *resmocks.AlertBroadcaster, *resmocks.PositionBroadcaster) {
	t.Helper()
	repo := repomocks.NewAlertRepository(t)
	broadcaster := resmocks.NewAlertBroadcaster(t)
	posBroadcaster := resmocks.NewPositionBroadcaster(t)
	c := controllers.NewAlert(logger.NewLogger(), repo, broadcaster, posBroadcaster)
	return c, repo, broadcaster, posBroadcaster
}

func TestAlert_List(t *testing.T) {
	t.Parallel()

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()
		c, repo, _, _ := newTestAlertController(t)

		alerts := []domain.Alert{testdata.GetTestAlert()}
		repo.On("List", mock.Anything, 20, 0).Return(alerts, int64(1), nil)

		req := httptest.NewRequest("GET", "/alerts?limit=20&offset=0", nil)
		rec := httptest.NewRecorder()

		c.List(rec, req)

		assert.Equal(t, 200, rec.Code)

		var response dtos.ListAlertsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, 20, response.Limit)
		assert.Equal(t, 0, response.Offset)
		assert.Len(t, response.Data, 1)
		assert.Equal(t, testdata.GetTestAlert().ID, response.Data[0].ID)
	})

	t.Run("works correctly when empty", func(t *testing.T) {
		t.Parallel()
		c, repo, _, _ := newTestAlertController(t)

		repo.On("List", mock.Anything, 20, 0).Return(nil, int64(0), nil)

		req := httptest.NewRequest("GET", "/alerts", nil)
		rec := httptest.NewRecorder()

		c.List(rec, req)

		assert.Equal(t, 200, rec.Code)

		var response dtos.ListAlertsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), response.Total)
		assert.Empty(t, response.Data)
	})

	t.Run("fails when repo fails", func(t *testing.T) {
		t.Parallel()
		c, repo, _, _ := newTestAlertController(t)

		repo.On("List", mock.Anything, 20, 0).Return(nil, int64(0), assert.AnError)

		req := httptest.NewRequest("GET", "/alerts", nil)
		rec := httptest.NewRecorder()

		c.List(rec, req)

		assert.Equal(t, 500, rec.Code)
	})
}

func TestAlert_Stream(t *testing.T) {
	t.Parallel()

	t.Run("works correctly and receives alert event", func(t *testing.T) {
		t.Parallel()

		c, _, broadcaster, posBroadcaster := newTestAlertController(t)
		alertCh := make(chan domain.Alert, 1)
		posCh := make(chan domain.Position, 1)

		broadcaster.On("Subscribe").Return((<-chan domain.Alert)(alertCh))
		broadcaster.On("Unsubscribe", mock.Anything).Return()
		posBroadcaster.On("Subscribe").Return((<-chan domain.Position)(posCh))
		posBroadcaster.On("Unsubscribe", mock.Anything).Return()

		req := httptest.NewRequest("GET", "/alerts/stream", nil)
		ctx, cancel := context.WithCancel(req.Context())
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		alert := testdata.GetTestAlert()
		alertCh <- alert

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		c.Stream(rec, req)

		assert.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
		assert.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
		assert.Equal(t, "keep-alive", rec.Header().Get("Connection"))

		body := rec.Body.String()
		assert.Contains(t, body, "event: alert")
		assert.Contains(t, body, "data:")

		dataLines := splitSSEData(body)
		if len(dataLines) > 0 {
			var event dtos.AlertSSEEvent
			err := json.Unmarshal(dataLines[0], &event)
			assert.NoError(t, err)
			assert.Equal(t, alert.ID, event.ID)
			assert.Equal(t, alert.VehicleID, event.VehicleID)
		}
	})

	t.Run("works correctly and receives position event", func(t *testing.T) {
		t.Parallel()

		c, _, broadcaster, posBroadcaster := newTestAlertController(t)
		alertCh := make(chan domain.Alert, 1)
		posCh := make(chan domain.Position, 1)

		broadcaster.On("Subscribe").Return((<-chan domain.Alert)(alertCh))
		broadcaster.On("Unsubscribe", mock.Anything).Return()
		posBroadcaster.On("Subscribe").Return((<-chan domain.Position)(posCh))
		posBroadcaster.On("Unsubscribe", mock.Anything).Return()

		req := httptest.NewRequest("GET", "/alerts/stream", nil)
		ctx, cancel := context.WithCancel(req.Context())
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		pos := testdata.GetTestPosition()
		posCh <- pos

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		c.Stream(rec, req)

		body := rec.Body.String()
		assert.Contains(t, body, "event: position")
		assert.Contains(t, body, "data:")

		dataLines := splitSSEData(body)
		if len(dataLines) > 0 {
			var event dtos.PositionSSEEvent
			err := json.Unmarshal(dataLines[0], &event)
			assert.NoError(t, err)
			assert.Equal(t, pos.VehicleID, event.VehicleID)
			assert.Equal(t, pos.Latitude, event.Latitude)
		}
	})

	t.Run("handles correctly when context cancelled", func(t *testing.T) {
		t.Parallel()

		c, _, broadcaster, posBroadcaster := newTestAlertController(t)
		alertCh := make(chan domain.Alert, 1)
		posCh := make(chan domain.Position, 1)

		broadcaster.On("Subscribe").Return((<-chan domain.Alert)(alertCh))
		broadcaster.On("Unsubscribe", mock.Anything).Return()
		posBroadcaster.On("Subscribe").Return((<-chan domain.Position)(posCh))
		posBroadcaster.On("Unsubscribe", mock.Anything).Return()

		req := httptest.NewRequest("GET", "/alerts/stream", nil)
		ctx, cancel := context.WithCancel(req.Context())
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		cancel()

		c.Stream(rec, req)

		assert.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	})
}

// splitSSEData extracts the JSON data payloads from an SSE response body.
func splitSSEData(body string) [][]byte {
	var result [][]byte
	current := ""
	for _, line := range splitLines(body) {
		if len(line) > 6 && line[:6] == "data: " {
			current = line[6:]
		} else if current != "" && line == "" {
			result = append(result, []byte(current))
			current = ""
		}
	}
	if current != "" {
		result = append(result, []byte(current))
	}
	return result
}

// splitLines splits a string by newlines without using strings.Split.
func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, c := range s {
		if c == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
