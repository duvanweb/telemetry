package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/core/ports/resources"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/dtos"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// Alert is the HTTP controller for alert-related endpoints.
type Alert struct {
	broadcaster resources.AlertBroadcaster
	logger      logger.Logger
}

// @Router /alerts/stream [get]
// @Tags alerts
// @Summary Stream alerts via Server-Sent Events.
// @Description Returns a text/event-stream with real-time alert notifications.
// @Success 200 {string} string "Event stream."
// Stream handles GET /alerts/stream and sends alert events via SSE.
func (c *Alert) Stream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.logger.Errorw(r.Context(), "streaming unsupported: response writer does not implement http.Flusher")
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	ch := c.broadcaster.Subscribe()
	defer c.broadcaster.Unsubscribe(ch)

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case alert, ok := <-ch:
			if !ok {
				return
			}
			c.writeSSEEvent(ctx, w, flusher, alert)
		}
	}
}

// writeSSEEvent writes a single SSE event to the response writer and flushes.
func (c *Alert) writeSSEEvent(ctx context.Context, w http.ResponseWriter, flusher http.Flusher, alert domain.Alert) {
	event := dtos.AlertSSEEvent{
		ID:         alert.ID,
		VehicleID:  alert.VehicleID,
		Type:       string(alert.Type),
		Latitude:   alert.Latitude,
		Longitude:  alert.Longitude,
		DetectedAt: alert.DetectedAt,
	}

	data, err := json.Marshal(event)
	if err != nil {
		c.logger.Errorw(ctx, "failed to marshal alert SSE event", "error", err)
		return
	}

	if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
		c.logger.Errorw(ctx, "failed to write SSE event", "error", err)
		return
	}

	flusher.Flush()
}

// NewAlert creates and returns a new Alert controller.
func NewAlert(log logger.Logger, broadcaster resources.AlertBroadcaster) *Alert {
	return &Alert{broadcaster: broadcaster, logger: log}
}
