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

func TestRepository_List(t *testing.T) {
	now := time.Now()

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := alertrepo.NewRepository(db)

		rows := sqlmock.NewRows([]string{"id", "vehicle_id", "type", "latitude", "longitude", "detected_at", "created_at"}).
			AddRow(int64(1), int64(1), domain.AlertTypeVehicleStopped, 4.71, -74.07, now, now).
			AddRow(int64(2), int64(2), domain.AlertTypeVehicleStopped, 4.65, -74.06, now.Add(-time.Hour), now.Add(-time.Hour))

		mock.ExpectQuery(alertrepo.ListAlertsQuery).
			WithArgs(20, 0).
			WillReturnRows(rows)
		mock.ExpectQuery(alertrepo.CountAlertsQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(2)))

		alerts, total, err := repo.List(context.Background(), 20, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, alerts, 2)
		assert.Equal(t, int64(1), alerts[0].ID)
		assert.Equal(t, int64(2), alerts[1].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("works correctly when empty", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := alertrepo.NewRepository(db)

		mock.ExpectQuery(alertrepo.ListAlertsQuery).
			WithArgs(20, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id", "vehicle_id", "type", "latitude", "longitude", "detected_at", "created_at"}))
		mock.ExpectQuery(alertrepo.CountAlertsQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))

		alerts, total, err := repo.List(context.Background(), 20, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, alerts)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("fails when query fails", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := alertrepo.NewRepository(db)

		mock.ExpectQuery(alertrepo.ListAlertsQuery).
			WithArgs(20, 0).
			WillReturnError(assert.AnError)

		alerts, total, err := repo.List(context.Background(), 20, 0)
		assert.Error(t, err)
		assert.Nil(t, alerts)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("fails when count fails", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := alertrepo.NewRepository(db)

		mock.ExpectQuery(alertrepo.ListAlertsQuery).
			WithArgs(20, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id", "vehicle_id", "type", "latitude", "longitude", "detected_at", "created_at"}))
		mock.ExpectQuery(alertrepo.CountAlertsQuery).
			WillReturnError(assert.AnError)

		alerts, total, err := repo.List(context.Background(), 20, 0)
		assert.Error(t, err)
		assert.Nil(t, alerts)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("fails when scan fails", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := alertrepo.NewRepository(db)

		// Wrong column type (string instead of int64 for id) causes scan error.
		rows := sqlmock.NewRows([]string{"id", "vehicle_id", "type", "latitude", "longitude", "detected_at", "created_at"}).
			AddRow("not-an-int", int64(1), domain.AlertTypeVehicleStopped, 4.71, -74.07, now, now)

		mock.ExpectQuery(alertrepo.ListAlertsQuery).
			WithArgs(20, 0).
			WillReturnRows(rows)

		alerts, total, err := repo.List(context.Background(), 20, 0)
		assert.Error(t, err)
		assert.Nil(t, alerts)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
