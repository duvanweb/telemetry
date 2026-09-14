package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

// Client is the HTTP client for calling vehicle-service.
type Client struct {
	httpClient *http.Client
	baseURL    string
	logger     logger.Logger
}

// ValidateVehicle checks if a vehicle exists by calling vehicle-service GET /api/vehicles/{id}.
// Returns nil if the vehicle exists, ErrVehicleNotFound if 404, ErrVehicleServiceUnavailable on 5xx/timeout.
func (c *Client) ValidateVehicle(ctx context.Context, vehicleID int64) error {
	url := fmt.Sprintf("%s/api/vehicles/%d", c.baseURL, vehicleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.logger.Errorw(ctx, "failed to create vehicle validation request", "error", err)
		return domain.ErrVehicleServiceUnavailable
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Errorw(ctx, "failed to call vehicle-service", "error", err)
		return domain.ErrVehicleServiceUnavailable
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return domain.ErrVehicleNotFound
	default:
		c.logger.Errorw(ctx, "vehicle-service returned unexpected status",
			"status", resp.StatusCode, "vehicle_id", vehicleID)
		return domain.ErrVehicleServiceUnavailable
	}
}

// NewVehicleClient creates and returns a new vehicle-service HTTP client.
func NewVehicleClient(config *env.Configuration, log logger.Logger) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    config.VehicleServiceURL,
		logger:     log,
	}
}
