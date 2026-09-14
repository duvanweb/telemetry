package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/telemetry-platform/vehicle-service/internal/core/domain"
	"github.com/telemetry-platform/vehicle-service/internal/core/vehicle"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/api/dtos"
	apierrors "github.com/telemetry-platform/vehicle-service/internal/infrastructure/api/errors"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
)

// Vehicle is the HTTP controller for vehicle-related endpoints.
type Vehicle struct {
	logger  logger.Logger
	service *vehicle.Service
}

// Create handles POST /api/vehicles requests.
//
// @Router /api/vehicles [post]
// @Tags vehicles
// @Summary Create a new vehicle.
// @Accept json
// @Produce json
// @Param request body dtos.CreateVehicleRequest true "Vehicle to create"
// @Success 201 {object} dtos.VehicleResponse "Vehicle created."
// @Failure 400 "Invalid plate format."
// @Failure 409 "Vehicle already exists."
func (c *Vehicle) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	v := &domain.Vehicle{Plate: req.Plate}
	if err := c.service.Create(r.Context(), v); err != nil {
		c.writeDomainError(w, r, err)
		return
	}

	response := toVehicleResponse(v)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

// FindByPlate handles GET /api/vehicles/plates/{plate} requests.
//
// @Router /api/vehicles/plates/{plate} [get]
// @Tags vehicles
// @Summary Get a vehicle by plate.
// @Produce json
// @Param plate path string true "Vehicle plate"
// @Success 200 {object} dtos.VehicleResponse "Vehicle found."
// @Failure 404 "Vehicle not found."
func (c *Vehicle) FindByPlate(w http.ResponseWriter, r *http.Request) {
	plate := chi.URLParam(r, "plate")

	v, err := c.service.FindByPlate(r.Context(), plate)
	if err != nil {
		c.writeDomainError(w, r, err)
		return
	}

	response := toVehicleResponse(v)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// GetByID handles GET /api/vehicles/{id} requests.
//
// @Router /api/vehicles/{id} [get]
// @Tags vehicles
// @Summary Get a vehicle by ID.
// @Produce json
// @Param id path int true "Vehicle ID"
// @Success 200 {object} dtos.VehicleResponse "Vehicle found."
// @Failure 404 "Vehicle not found."
func (c *Vehicle) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	v, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		c.writeDomainError(w, r, err)
		return
	}

	response := toVehicleResponse(v)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// List handles GET /api/vehicles requests.
//
// @Router /api/vehicles [get]
// @Tags vehicles
// @Summary List active vehicles with pagination.
// @Produce json
// @Param limit query int false "Page limit (default 20, max 100)"
// @Param offset query int false "Page offset (default 0)"
// @Success 200 {object} dtos.ListVehiclesResponse "Paginated list of vehicles."
func (c *Vehicle) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)

	vehicles, total, err := c.service.List(r.Context(), limit, offset)
	if err != nil {
		c.writeDomainError(w, r, err)
		return
	}

	data := make([]dtos.VehicleResponse, 0, len(vehicles))
	for i := range vehicles {
		data = append(data, toVehicleResponse(&vehicles[i]))
	}

	response := dtos.ListVehiclesResponse{
		Data:   data,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// NewVehicle creates and returns a new Vehicle controller.
func NewVehicle(log logger.Logger, svc *vehicle.Service) *Vehicle {
	return &Vehicle{logger: log, service: svc}
}

// SoftDelete handles DELETE /api/vehicles/{id} requests.
//
// @Router /api/vehicles/{id} [delete]
// @Tags vehicles
// @Summary Soft-delete a vehicle.
// @Param id path int true "Vehicle ID"
// @Success 204 "Vehicle soft-deleted."
// @Failure 404 "Vehicle not found."
func (c *Vehicle) SoftDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := c.service.SoftDelete(r.Context(), id); err != nil {
		c.writeDomainError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// parseID extracts and parses the vehicle ID from the URL path.
func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

// parsePagination extracts and normalizes limit/offset from query parameters.
// Applies defaults limit=20, offset=0 and caps limit at 100.
func parsePagination(r *http.Request) (limit, offset int) {
	limit = 20
	offset = 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	if limit > 100 {
		limit = 100
	}

	return limit, offset
}

// toVehicleResponse converts a domain Vehicle to a VehicleResponse DTO.
func toVehicleResponse(v *domain.Vehicle) dtos.VehicleResponse {
	return dtos.VehicleResponse{
		ID:        v.ID,
		Plate:     v.Plate,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

// writeDomainError maps a domain error to the appropriate HTTP status code.
func (c *Vehicle) writeDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidPlate):
		apierrors.WriteError(w, http.StatusBadRequest, err)
	case errors.Is(err, domain.ErrVehicleNotFound):
		apierrors.WriteError(w, http.StatusNotFound, err)
	case errors.Is(err, domain.ErrVehicleAlreadyExists):
		apierrors.WriteError(w, http.StatusConflict, err)
	default:
		c.logger.Errorw(r.Context(), "unexpected vehicle error", "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
	}
}
