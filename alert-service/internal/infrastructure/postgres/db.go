package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

// NewDB creates and returns a new PostgreSQL database connection.
// It pings the database to verify connectivity before returning.
func NewDB(config *env.Configuration, log logger.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Infow(context.Background(), "database connection established",
		"url", config.PostgresURL)

	return db, nil
}
