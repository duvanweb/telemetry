package position_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	positionrepo "github.com/telemetry-platform/geo-service/internal/infrastructure/postgres/repositories/position"
	testdata "github.com/telemetry-platform/geo-service/test/data"
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
	t.Parallel()

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := positionrepo.NewRepository(db)
		pos := testdata.GetTestPosition()

		mock.ExpectExec(positionrepo.SavePositionQuery).
			WithArgs(pos.VehicleID, pos.Latitude, pos.Longitude, pos.RecordedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(context.Background(), pos)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("fails when query fails", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := positionrepo.NewRepository(db)
		pos := testdata.GetTestPosition()

		mock.ExpectExec(positionrepo.SavePositionQuery).
			WithArgs(pos.VehicleID, pos.Latitude, pos.Longitude, pos.RecordedAt).
			WillReturnError(assert.AnError)

		err := repo.Save(context.Background(), pos)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_DeleteByVehicleID(t *testing.T) {
	t.Parallel()

	t.Run("works correctly", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := positionrepo.NewRepository(db)

		mock.ExpectExec(positionrepo.DeleteByVehicleIDQuery).
			WithArgs(int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 5))

		err := repo.DeleteByVehicleID(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("works correctly when no positions exist", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := positionrepo.NewRepository(db)

		mock.ExpectExec(positionrepo.DeleteByVehicleIDQuery).
			WithArgs(int64(999)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteByVehicleID(context.Background(), 999)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("fails when query fails", func(t *testing.T) {
		t.Parallel()
		db, mock := newMockDB(t)
		repo := positionrepo.NewRepository(db)

		mock.ExpectExec(positionrepo.DeleteByVehicleIDQuery).
			WithArgs(int64(1)).
			WillReturnError(assert.AnError)

		err := repo.DeleteByVehicleID(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
