package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	"github.com/telemetry-platform/geo-service/internal/core/ports/services"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/api/dtos"
	apierrors "github.com/telemetry-platform/geo-service/internal/infrastructure/api/errors"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// Position is the HTTP controller for position-related endpoints.
type Position struct {
	logger  logger.Logger
	service services.PositionService
}

// @Router /api/vehicles/{vehicle_id}/positions [post]
// @Tags positions
// @Summary Ingest a GPS position for a vehicle.
// @Accept json
// @Produce json
// @Param vehicle_id path int true "Vehicle ID"
// @Param request body dtos.CreatePositionRequest true "Position data"
// @Success 202 {object} dtos.CreatePositionResponse "Position accepted for processing."
// @Failure 400 "Invalid request body or coordinates."
// @Failure 404 "Vehicle not found."
// @Failure 409 "Duplicate position within TTL window."
// @Failure 503 "Vehicle service unavailable."
// Create handles POST /api/vehicles/{vehicle_id}/positions requests.
func (c *Position) Create(w http.ResponseWriter, r *http.Request) {
	vehicleID, err := strconv.ParseInt(chi.URLParam(r, "vehicle_id"), 10, 64)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var req dtos.CreatePositionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	recordedAt, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	pos := domain.Position{
		VehicleID:  vehicleID,
		Latitude:   req.Lat,
		Longitude:  req.Lng,
		RecordedAt: recordedAt,
	}

	if err := c.service.Ingest(r.Context(), pos); err != nil {
		c.writeDomainError(w, r, err)
		return
	}

	response := dtos.CreatePositionResponse{Status: "accepted", VehicleID: vehicleID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(response)
}

// NewPosition creates and returns a new Position controller.
func NewPosition(log logger.Logger, svc services.PositionService) *Position {
	return &Position{logger: log, service: svc}
}

// writeDomainError maps a domain error to the appropriate HTTP status code.
func (c *Position) writeDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidPosition):
		apierrors.WriteError(w, http.StatusBadRequest, err)
	case errors.Is(err, domain.ErrVehicleNotFound):
		apierrors.WriteError(w, http.StatusNotFound, err)
	case errors.Is(err, domain.ErrDuplicatePosition):
		apierrors.WriteError(w, http.StatusConflict, err)
	case errors.Is(err, domain.ErrVehicleServiceUnavailable):
		apierrors.WriteError(w, http.StatusServiceUnavailable, err)
	default:
		c.logger.Errorw(r.Context(), "unexpected position error", "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
	}
}
