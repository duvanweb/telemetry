package position_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetry-platform/geo-service/internal/core/domain"
	positionrepo "github.com/telemetry-platform/geo-service/internal/infrastructure/postgres/repositories/position"
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
	pos := domain.Position{VehicleID: 1, Latitude: 4.71, Longitude: -74.07, RecordedAt: now}

	t.Run("works correctly", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := positionrepo.NewRepository(db)

		mock.ExpectExec(positionrepo.SavePositionQuery).
			WithArgs(int64(1), 4.71, -74.07, now).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(context.Background(), pos)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("fails when query fails", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := positionrepo.NewRepository(db)

		mock.ExpectExec(positionrepo.SavePositionQuery).
			WithArgs(int64(1), 4.71, -74.07, now).
			WillReturnError(assert.AnError)

		err := repo.Save(context.Background(), pos)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
