package domain

import "errors"

var (
	ErrAlertAlreadyExists = errors.New("alert already exists for this stop event")
	ErrInvalidAlertType   = errors.New("invalid alert type")
	ErrTrackNotFound      = errors.New("vehicle track not found")
)
