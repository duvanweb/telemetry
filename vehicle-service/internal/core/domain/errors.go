package domain

import "errors"

var (
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrVehicleAlreadyExists = errors.New("vehicle already exists")
	ErrInvalidPlate         = errors.New("invalid plate format, expected ABC-123")
)
