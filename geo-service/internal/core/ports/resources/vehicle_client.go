package resources

import "context"

// VehicleClient is the port for validating vehicles against vehicle-service.
//
//go:generate mockery --name VehicleClient --dir=. --output=./mocks
type VehicleClient interface {
	ValidateVehicle(ctx context.Context, vehicleID int64) error
}
