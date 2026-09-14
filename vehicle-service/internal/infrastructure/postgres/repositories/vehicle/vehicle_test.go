package vehicle

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetry-platform/vehicle-service/internal/core/domain"
	testdata "github.com/telemetry-platform/vehicle-service/test/data"
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

func TestRepository_Create(t *testing.T) {
	t.Run("works correctly", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		mock.ExpectQuery(createVehicleQuery).
			WithArgs("ABC-123").
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(int64(1), now, now))

		v := &domain.Vehicle{Plate: testdata.GetTestVehicle().Plate}
		err := repo.Create(context.Background(), v)

		assert.NoError(t, err)
		assert.Equal(t, testdata.GetTestVehicleID(), v.ID)
		assert.NotZero(t, v.CreatedAt)
		assert.NotZero(t, v.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("handles correctly when unique constraint violated", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		mock.ExpectQuery(createVehicleQuery).
			WithArgs("ABC-123").
			WillReturnError(&pgconn.PgError{
				Code:    uniqueViolationCode,
				Message: "duplicate key value violates unique constraint",
			})

		v := &domain.Vehicle{Plate: testdata.GetTestVehicle().Plate}
		err := repo.Create(context.Background(), v)

		assert.Equal(t, domain.ErrVehicleAlreadyExists, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_FindByPlate(t *testing.T) {
	t.Run("works correctly", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		expected := testdata.GetTestVehicle()
		mock.ExpectQuery(findVehicleByPlateQuery).
			WithArgs(expected.Plate).
			WillReturnRows(sqlmock.NewRows([]string{"id", "plate", "created_at", "updated_at"}).
				AddRow(expected.ID, expected.Plate, expected.CreatedAt, expected.UpdatedAt))

		v, err := repo.FindByPlate(context.Background(), expected.Plate)

		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, expected.ID, v.ID)
		assert.Equal(t, expected.Plate, v.Plate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("handles correctly when vehicle not found", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		mock.ExpectQuery(findVehicleByPlateQuery).
			WithArgs("ZZZ-999").
			WillReturnError(sql.ErrNoRows)

		v, err := repo.FindByPlate(context.Background(), "ZZZ-999")

		assert.Equal(t, domain.ErrVehicleNotFound, err)
		assert.Nil(t, v)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_GetByID(t *testing.T) {
	t.Run("works correctly", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		expected := testdata.GetTestVehicle()
		mock.ExpectQuery(getVehicleByIDQuery).
			WithArgs(expected.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "plate", "created_at", "updated_at"}).
				AddRow(expected.ID, expected.Plate, expected.CreatedAt, expected.UpdatedAt))

		v, err := repo.GetByID(context.Background(), expected.ID)

		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, expected.ID, v.ID)
		assert.Equal(t, expected.Plate, v.Plate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("handles correctly when vehicle not found", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		mock.ExpectQuery(getVehicleByIDQuery).
			WithArgs(int64(999)).
			WillReturnError(sql.ErrNoRows)

		v, err := repo.GetByID(context.Background(), 999)

		assert.Equal(t, domain.ErrVehicleNotFound, err)
		assert.Nil(t, v)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_List(t *testing.T) {
	t.Run("works correctly", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		expected := testdata.GetTestVehicles()
		mock.ExpectQuery(listVehiclesQuery).
			WithArgs(20, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id", "plate", "created_at", "updated_at"}).
				AddRow(expected[0].ID, expected[0].Plate, expected[0].CreatedAt, expected[0].UpdatedAt).
				AddRow(expected[1].ID, expected[1].Plate, expected[1].CreatedAt, expected[1].UpdatedAt))

		mock.ExpectQuery(countVehiclesQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(len(expected))))

		vehicles, total, err := repo.List(context.Background(), 20, 0)

		assert.NoError(t, err)
		assert.Len(t, vehicles, len(expected))
		assert.Equal(t, int64(len(expected)), total)
		assert.Equal(t, expected[0].Plate, vehicles[0].Plate)
		assert.Equal(t, expected[1].Plate, vehicles[1].Plate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_SoftDelete(t *testing.T) {
	t.Run("works correctly", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		mock.ExpectExec(softDeleteVehicleQuery).
			WithArgs(testdata.GetTestVehicleID()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.SoftDelete(context.Background(), testdata.GetTestVehicleID())

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("handles correctly when vehicle not found", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		mock.ExpectExec(softDeleteVehicleQuery).
			WithArgs(int64(999)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.SoftDelete(context.Background(), 999)

		assert.Equal(t, domain.ErrVehicleNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
