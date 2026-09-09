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
	resmocks "github.com/telemetry-platform/alert-service/internal/core/ports/resources/mocks"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/controllers"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/dtos"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
	testdata "github.com/telemetry-platform/alert-service/test/data"
)

func TestAlert_Stream(t *testing.T) {
	t.Parallel()

	t.Run("works correctly and receives events", func(t *testing.T) {
		t.Parallel()

		broadcaster := resmocks.NewAlertBroadcaster(t)
		ch := make(chan domain.Alert, 1)

		broadcaster.On("Subscribe").Return((<-chan domain.Alert)(ch))
		broadcaster.On("Unsubscribe", mock.Anything).Return()

		c := controllers.NewAlert(logger.NewLogger(), broadcaster)

		req := httptest.NewRequest("GET", "/alerts/stream", nil)
		ctx, cancel := context.WithCancel(req.Context())
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		// Send an alert then cancel the context to stop the stream.
		alert := testdata.GetTestAlert()
		ch <- alert

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		c.Stream(rec, req)

		assert.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
		assert.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
		assert.Equal(t, "keep-alive", rec.Header().Get("Connection"))

		// Verify the SSE event was written.
		body := rec.Body.String()
		assert.Contains(t, body, "data:")

		// Parse the SSE data line to verify the alert content.
		lines := splitSSEData(body)
		if len(lines) > 0 {
			var event dtos.AlertSSEEvent
			err := json.Unmarshal(lines[0], &event)
			assert.NoError(t, err)
			assert.Equal(t, alert.ID, event.ID)
			assert.Equal(t, alert.VehicleID, event.VehicleID)
		}
	})

	t.Run("handles correctly when context cancelled", func(t *testing.T) {
		t.Parallel()

		broadcaster := resmocks.NewAlertBroadcaster(t)
		ch := make(chan domain.Alert, 1)

		broadcaster.On("Subscribe").Return((<-chan domain.Alert)(ch))
		broadcaster.On("Unsubscribe", mock.Anything).Return()

		c := controllers.NewAlert(logger.NewLogger(), broadcaster)

		req := httptest.NewRequest("GET", "/alerts/stream", nil)
		ctx, cancel := context.WithCancel(req.Context())
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		// Cancel immediately.
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
