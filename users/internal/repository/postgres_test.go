package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func TestIntegrationInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbConn := setupDB(t)
	defer teardownDB(t, dbConn)

	repo := NewPostgres(dbConn)

	user := &User{
		ID:           uuid.New().String(),
		Fullname:     "John Doe",
		Username:     "jdoe",
		Birthdate:    "2000-01-01",
		Email:        "joedoe@mail.com",
		PasswordHash: "123456",
		Role:         "user",
		CreatedAt:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	actual, err := repo.Insert(context.TODO(), user)
	require.NoError(t, err)

	require.Equal(t, user, actual)
}

func TestIntegrationExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbConn := setupDB(t)
	defer teardownDB(t, dbConn)

	repo := NewPostgres(dbConn)

	user := &User{
		ID:           uuid.New().String(),
		Fullname:     "John Doe",
		Username:     "jdoe",
		Birthdate:    "2000-01-01",
		Email:        "",
		PasswordHash: "123456",
		Role:         "user",
		CreatedAt:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	_, err := repo.Insert(context.TODO(), user)
	require.NoError(t, err)

	t.Run("user exists", func(t *testing.T) {
		actual, err := repo.Exists(context.TODO(), user.Username, user.Email)
		require.NoError(t, err)

		require.True(t, actual)
	})

	t.Run("user does not exist", func(t *testing.T) {
		actual, err := repo.Exists(context.TODO(), "foo", "foo@bar.quz")
		require.NoError(t, err)

		require.False(t, actual)
	})
}

func TestIntegrationSelectByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbConn := setupDB(t)
	defer teardownDB(t, dbConn)

	repo := NewPostgres(dbConn)

	user := &User{
		ID:           uuid.New().String(),
		Fullname:     "John Doe",
		Username:     "jdoe",
		Birthdate:    "2000-01-01",
		Email:        "joedoe@mail.com",
		PasswordHash: "123456",
		Role:         "user",
		CreatedAt:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	_, err := repo.Insert(context.TODO(), user)
	require.NoError(t, err)

	t.Run("user exists", func(t *testing.T) {
		actual, err := repo.SelectByID(context.TODO(), user.ID)
		require.NoError(t, err)

		require.Equal(t, user, actual)
	})

	t.Run("user does not exist", func(t *testing.T) {
		actual, err := repo.SelectByID(context.TODO(), uuid.New().String())
		require.NoError(t, err)

		require.Nil(t, actual)
	})
}

func TestIntegrationSelectByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbConn := setupDB(t)
	defer teardownDB(t, dbConn)

	repo := NewPostgres(dbConn)

	user := &User{
		ID:           uuid.New().String(),
		Fullname:     "John Doe",
		Username:     "jdoe",
		Birthdate:    "2000-01-01",
		Email:        "joedoe@mail.com",
		PasswordHash: "123456",
		Role:         "user",
		CreatedAt:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	_, err := repo.Insert(context.TODO(), user)
	require.NoError(t, err)

	t.Run("user exists", func(t *testing.T) {
		actual, err := repo.SelectByEmail(context.TODO(), user.Email)
		require.NoError(t, err)

		require.Equal(t, user, actual)
	})

	t.Run("user does not exist", func(t *testing.T) {
		actual, err := repo.SelectByEmail(context.TODO(), "foo@bar.quz")
		require.NoError(t, err)

		require.Nil(t, actual)
	})
}

func setupDB(t *testing.T) *sqlx.DB {
	dbConn, err := sqlx.Connect("pgx", "postgres://user:password@localhost:5432/testdb?sslmode=disable")
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		err = dbConn.Ping()
		if err == nil {
			break
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
	require.NoError(t, err)

	return dbConn
}

func teardownDB(t *testing.T, dbConn *sqlx.DB) {
	_, err := dbConn.Exec("TRUNCATE TABLE users CASCADE")
	require.NoError(t, err)

	require.NoError(t, dbConn.Close())
}
