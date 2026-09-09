package alert_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetry-platform/alert-service/internal/core/domain"
	alertrepo "github.com/telemetry-platform/alert-service/internal/infrastructure/postgres/repositories/alert"
)

// newMockDB creates a sqlmock DB with exact query matching for testing.
func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})
	return db, mock
}

func TestRepository_Save(t *testing.T) {
	now := time.Now()
	alert := domain.Alert{
		VehicleID:  1,
		Type:       domain.AlertTypeVehicleStopped,
		Latitude:   4.71,
		Longitude:  -74.07,
		DetectedAt: now,
	}

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := alertrepo.NewRepository(db)

		mock.ExpectQuery(alertrepo.SaveAlertQuery).
			WithArgs(int64(1), domain.AlertTypeVehicleStopped, 4.71, -74.07, now).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(int64(1), now))

		err := repo.Save(context.Background(), &alert)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), alert.ID)
		assert.Equal(t, now, alert.CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("fails when query fails", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := alertrepo.NewRepository(db)

		mock.ExpectQuery(alertrepo.SaveAlertQuery).
			WithArgs(int64(1), domain.AlertTypeVehicleStopped, 4.71, -74.07, now).
			WillReturnError(assert.AnError)

		err := repo.Save(context.Background(), &alert)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
