package domain

import "errors"

var (
	ErrVehicleNotFound           = errors.New("vehicle not found")
	ErrVehicleServiceUnavailable = errors.New("vehicle service unavailable")
	ErrDuplicatePosition         = errors.New("duplicate position")
	ErrInvalidPosition           = errors.New("invalid position, lat or lng out of range")
)
