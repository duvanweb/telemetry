package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/telemetry-platform/vehicle-service/internal/core/ports/repositories"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/env"
	"github.com/telemetry-platform/vehicle-service/internal/infrastructure/pkg/logger"
)

// NewConnection creates and returns a new PostgreSQL database connection as a Databaser.
// It configures the connection pool and pings the database to verify connectivity.
// If the ping fails, it retries with exponential backoff (5 attempts: 2s, 4s, 8s, 16s)
// before giving up, so transient DB unavailability at startup does not crash the service.
func NewConnection(config *env.Configuration, log logger.Logger) (repositories.Databaser, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	const maxRetries = 5
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			break
		}
		log.Warnw(context.Background(), "database ping failed, retrying",
			"attempt", attempt, "of", maxRetries, "error", err)
		if attempt < maxRetries {
			time.Sleep(time.Duration(1<<attempt) * time.Second) // 2s, 4s, 8s, 16s
		}
	}
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, err)
	}

	log.Infow(context.Background(), "database connection established",
		"host", config.DBHost, "port", config.DBPort, "db", config.DBName)

	return db, nil
}
