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
	t.Run("success", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		now := time.Now()
		mock.ExpectQuery(createVehicleQuery).
			WithArgs("ABC-123").
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(int64(1), now, now))

		v := &domain.Vehicle{Plate: "ABC-123"}
		err := repo.Create(context.Background(), v)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), v.ID)
		assert.NotZero(t, v.CreatedAt)
		assert.NotZero(t, v.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("unique violation", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		mock.ExpectQuery(createVehicleQuery).
			WithArgs("ABC-123").
			WillReturnError(&pgconn.PgError{
				Code:    uniqueViolationCode,
				Message: "duplicate key value violates unique constraint",
			})

		v := &domain.Vehicle{Plate: "ABC-123"}
		err := repo.Create(context.Background(), v)

		assert.Equal(t, domain.ErrVehicleAlreadyExists, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_SoftDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		mock.ExpectExec(softDeleteVehicleQuery).
			WithArgs(int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.SoftDelete(context.Background(), 1)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
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

func TestRepository_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		now := time.Now()
		mock.ExpectQuery(getVehicleByIDQuery).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "plate", "created_at", "updated_at"}).
				AddRow(int64(1), "ABC-123", now, now))

		v, err := repo.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, int64(1), v.ID)
		assert.Equal(t, "ABC-123", v.Plate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
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

func TestRepository_FindByPlate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewRepository(db)

		now := time.Now()
		mock.ExpectQuery(findVehicleByPlateQuery).
			WithArgs("ABC-123").
			WillReturnRows(sqlmock.NewRows([]string{"id", "plate", "created_at", "updated_at"}).
				AddRow(int64(1), "ABC-123", now, now))

		v, err := repo.FindByPlate(context.Background(), "ABC-123")

		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, int64(1), v.ID)
		assert.Equal(t, "ABC-123", v.Plate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
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

func TestRepository_List(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewRepository(db)

	now := time.Now()
	mock.ExpectQuery(listVehiclesQuery).
		WithArgs(20, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "plate", "created_at", "updated_at"}).
			AddRow(int64(1), "ABC-123", now, now).
			AddRow(int64(2), "DEF-456", now, now))

	mock.ExpectQuery(countVehiclesQuery).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(2)))

	vehicles, total, err := repo.List(context.Background(), 20, 0)

	assert.NoError(t, err)
	assert.Len(t, vehicles, 2)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, "ABC-123", vehicles[0].Plate)
	assert.Equal(t, "DEF-456", vehicles[1].Plate)
	assert.NoError(t, mock.ExpectationsWereMet())
}
