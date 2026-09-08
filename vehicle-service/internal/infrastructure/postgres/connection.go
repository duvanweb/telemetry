package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
)

// NewDB creates and returns a new PostgreSQL database connection.
// It pings the database to verify connectivity before returning.
func NewDB(config *env.Configuration, log logger.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Infow(context.Background(), "database connection established",
		"host", config.DBHost, "port", config.DBPort, "db", config.DBName)

	return db, nil
}
