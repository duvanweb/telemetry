package domain

import "time"

// Vehicle is the core vehicle entity.
type Vehicle struct {
	ID        int64
	Plate     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time // nil = active (not soft-deleted)
}
