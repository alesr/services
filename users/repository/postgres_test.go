package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const dbConnStr string = "postgres://user:password@localhost:5432/testdb?sslmode=disable"

func TestIntegrationInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("user is inserted", func(t *testing.T) {
		dbConn := setupDB(t)
		defer teardownDB(t, dbConn)

		repo := NewPostgres(dbConn)

		user := &User{
			ID:            uuid.New().String(),
			Fullname:      "John Doe",
			Username:      "jdoe",
			Birthdate:     "2000-01-01",
			Email:         "joedoe@mail.com",
			EmailVerified: false,
			PasswordHash:  "123456",
			Role:          "user",
			CreatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		actual, err := repo.Insert(context.TODO(), user)
		require.NoError(t, err)

		require.Equal(t, user, actual)
	})

	t.Run("cannot insert the same user twice", func(t *testing.T) {
		dbConn := setupDB(t)
		defer teardownDB(t, dbConn)

		repo := NewPostgres(dbConn)

		user := &User{
			ID:            uuid.New().String(),
			Fullname:      "John Doe",
			Username:      "jdoe",
			Birthdate:     "2000-01-01",
			Email:         "joedoe@mail.com",
			EmailVerified: false,
			PasswordHash:  "123456",
			Role:          "user",
			CreatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		_, err := repo.Insert(context.TODO(), user)
		require.NoError(t, err)

		_, err = repo.Insert(context.TODO(), user)
		assert.Error(t, err)
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
		ID:            uuid.New().String(),
		Fullname:      "John Doe",
		Username:      "jdoe",
		Birthdate:     "2000-01-01",
		Email:         "joedoe@mail.com",
		EmailVerified: false,
		PasswordHash:  "123456",
		Role:          "user",
		CreatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
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

	t.Run("user is deleted", func(t *testing.T) {
		err := repo.DeleteByID(context.TODO(), user.ID)
		require.NoError(t, err)

		actual, err := repo.SelectByID(context.TODO(), user.ID)
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
		ID:            uuid.New().String(),
		Fullname:      "John Doe",
		Username:      "jdoe",
		Birthdate:     "2000-01-01",
		Email:         "joedoe@mail.com",
		EmailVerified: false,
		PasswordHash:  "123456",
		Role:          "user",
		CreatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
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

	t.Run("user is deleted", func(t *testing.T) {
		err := repo.DeleteByID(context.TODO(), user.ID)
		require.NoError(t, err)

		actual, err := repo.SelectByID(context.TODO(), user.ID)
		require.NoError(t, err)

		require.Nil(t, actual)
	})
}

func TestIntegrationDeleteByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbConn := setupDB(t)
	defer teardownDB(t, dbConn)

	repo := NewPostgres(dbConn)

	user := &User{
		ID:            uuid.New().String(),
		Fullname:      "John Doe",
		Username:      "jdoe",
		Birthdate:     "2000-01-01",
		Email:         "joedoe@mail.com",
		EmailVerified: false,
		PasswordHash:  "123456",
		Role:          "user",
		CreatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	_, err := repo.Insert(context.TODO(), user)
	require.NoError(t, err)

	t.Run("user exists", func(t *testing.T) {
		err := repo.DeleteByID(context.TODO(), user.ID)
		require.NoError(t, err)

		actual, err := repo.SelectByID(context.TODO(), user.ID)
		require.NoError(t, err)

		require.Nil(t, actual)
	})

	t.Run("user does not exist", func(t *testing.T) {
		err := repo.DeleteByID(context.TODO(), uuid.New().String())
		assert.Equal(t, ErrRecordNotFound, err)
	})
}

func TestIntegrationInsertEmailVerification(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbConn := setupDB(t)
	defer teardownDB(t, dbConn)

	repo := NewPostgres(dbConn)

	userID := uuid.New().String()
	user := &User{
		ID:            userID,
		Fullname:      "John Doe",
		Username:      "jdoe",
		Birthdate:     "2000-01-01",
		Email:         "joedoe@mail.com",
		EmailVerified: false,
		PasswordHash:  "123456",
		Role:          "user",
		CreatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	_, err := repo.Insert(context.TODO(), user)
	require.NoError(t, err)

	emailVerification := EmailVerification{
		Code:      "123456",
		UserID:    userID,
		CreatedAt: time.Time{},
		ExpiresAt: time.Time{},
	}

	err = repo.InsertEmailVerification(context.TODO(), emailVerification)
	require.NoError(t, err)
}

func setupDB(t *testing.T) *sqlx.DB {
	dbConn, err := sqlx.Connect("pgx", dbConnStr)
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

	_, err = dbConn.Exec("TRUNCATE TABLE email_verifications CASCADE")
	require.NoError(t, err)

	require.NoError(t, dbConn.Close())
}
