package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	"github.com/telemetry-platform/alert-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/alert-service/internal/core/ports/resources"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/api/dtos"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// Alert is the HTTP controller for alert-related endpoints.
type Alert struct {
	repo            repositories.AlertRepository
	broadcaster     resources.AlertBroadcaster
	posBroadcaster  resources.PositionBroadcaster
	logger          logger.Logger
}

// List handles GET /alerts and returns a paginated list of alerts.
//
// @Router /alerts [get]
// @Tags alerts
// @Summary List alerts with pagination.
// @Description Returns a paginated list of alerts ordered by detected_at descending.
// @Param limit query int false "Page size (default 20, max 100)."
// @Param offset query int false "Page offset (default 0)."
// @Success 200 {object} dtos.ListAlertsResponse "Paginated list of alerts."
// @Failure 500 "Unexpected error."
func (c *Alert) List(w http.ResponseWriter, r *http.Request) {
	limit := parseQueryInt(r, "limit", 20, 100)
	offset := parseQueryInt(r, "offset", 0, 0)

	alerts, total, err := c.repo.List(r.Context(), limit, offset)
	if err != nil {
		c.logger.Errorw(r.Context(), "failed to list alerts", "error", err)
		http.Error(w, "failed to list alerts", http.StatusInternalServerError)
		return
	}

	response := dtos.ListAlertsResponse{
		Data:   make([]dtos.AlertResponse, 0, len(alerts)),
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
	for _, alert := range alerts {
		response.Data = append(response.Data, toAlertResponse(alert))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(response); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode response", "error", encErr)
	}
}

// NewAlert creates and returns a new Alert controller.
func NewAlert(log logger.Logger, repo repositories.AlertRepository, broadcaster resources.AlertBroadcaster, posBroadcaster resources.PositionBroadcaster) *Alert {
	return &Alert{repo: repo, broadcaster: broadcaster, posBroadcaster: posBroadcaster, logger: log}
}

// Stream handles GET /alerts/stream and sends alert and position events via SSE.
//
// @Router /alerts/stream [get]
// @Tags alerts
// @Summary Stream alerts and positions via Server-Sent Events.
// @Description Returns a text/event-stream with real-time alert and position events.
// @Success 200 {string} string "Event stream."
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

	alertCh := c.broadcaster.Subscribe()
	defer c.broadcaster.Unsubscribe(alertCh)

	posCh := c.posBroadcaster.Subscribe()
	defer c.posBroadcaster.Unsubscribe(posCh)

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case alert, ok := <-alertCh:
			if !ok {
				return
			}
			c.writeSSEEvent(ctx, w, flusher, "alert", alert)
		case pos, ok := <-posCh:
			if !ok {
				return
			}
			c.writePositionSSEEvent(ctx, w, flusher, pos)
		}
	}
}

// writePositionSSEEvent writes a single position SSE event to the response writer and flushes.
func (c *Alert) writePositionSSEEvent(ctx context.Context, w http.ResponseWriter, flusher http.Flusher, pos domain.Position) {
	event := dtos.PositionSSEEvent{
		VehicleID:  pos.VehicleID,
		Latitude:   pos.Latitude,
		Longitude:  pos.Longitude,
		RecordedAt: pos.RecordedAt,
	}

	data, err := json.Marshal(event)
	if err != nil {
		c.logger.Errorw(ctx, "failed to marshal position SSE event", "error", err)
		return
	}

	if _, err := fmt.Fprintf(w, "event: position\ndata: %s\n\n", data); err != nil {
		c.logger.Errorw(ctx, "failed to write SSE event", "error", err)
		return
	}

	flusher.Flush()
}

// writeSSEEvent writes a single alert SSE event to the response writer and flushes.
func (c *Alert) writeSSEEvent(ctx context.Context, w http.ResponseWriter, flusher http.Flusher, eventType string, alert domain.Alert) {
	event := dtos.AlertSSEEvent{
		ID:         alert.ID,
		VehicleID:  alert.VehicleID,
		Type:       string(alert.Type),
		Latitude:   alert.Latitude,
		Longitude:  alert.Longitude,
		DetectedAt: alert.DetectedAt,
		CreatedAt:  alert.CreatedAt,
	}

	data, err := json.Marshal(event)
	if err != nil {
		c.logger.Errorw(ctx, "failed to marshal alert SSE event", "error", err)
		return
	}

	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, data); err != nil {
		c.logger.Errorw(ctx, "failed to write SSE event", "error", err)
		return
	}

	flusher.Flush()
}

// parseQueryInt parses an integer query parameter with a default and maximum value.
func parseQueryInt(r *http.Request, key string, defaultValue, maxValue int) int {
	value := defaultValue
	if raw := r.URL.Query().Get(key); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			value = parsed
		}
	}
	if maxValue > 0 && value > maxValue {
		value = maxValue
	}
	return value
}

// toAlertResponse maps a domain alert to an AlertResponse DTO.
func toAlertResponse(alert domain.Alert) dtos.AlertResponse {
	return dtos.AlertResponse{
		ID:         alert.ID,
		VehicleID:  alert.VehicleID,
		Type:       string(alert.Type),
		Latitude:   alert.Latitude,
		Longitude:  alert.Longitude,
		DetectedAt: alert.DetectedAt,
		CreatedAt:  alert.CreatedAt,
	}
}
