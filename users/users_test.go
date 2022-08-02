package users

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alesr/services/users/internal/repository"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	givenJWTSigningKey := "secret"
	givenRepo := &repositoryMock{}

	actual := New(givenJWTSigningKey, givenRepo)

	require.NotNil(t, actual)
	assert.Equal(t, givenJWTSigningKey, actual.jwtSigningKey)
	assert.Equal(t, givenRepo, actual.repo)
}

func TestCreate_validation(t *testing.T) {
	t.Parallel()

	given := CreateUserInput{
		Fullname: "%invalid-name%",
	}

	svc := Service{}

	_, err := svc.Create(context.Background(), given)
	assert.Error(t, err)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	givenUser := CreateUserInput{
		Fullname:        "John Doe",
		Username:        "jdoe",
		Birthdate:       "2000-01-01",
		Email:           "joedoe@mail.com",
		Password:        "password#123",
		ConfirmPassword: "password#123",
		Role:            RoleUser,
	}

	testCases := []struct {
		name          string
		givenUser     CreateUserInput
		givenRepoMock *repositoryMock
		expectedUser  *User
		expectedError error
	}{
		{
			name:      "user aleready exists",
			givenUser: givenUser,
			givenRepoMock: &repositoryMock{
				existsFunc: func(ctx context.Context, username, email string) (bool, error) {
					return true, nil
				},
			},
			expectedUser:  nil,
			expectedError: errAlreadyExists,
		},
		{
			name:      "check if user already exists error",
			givenUser: givenUser,
			givenRepoMock: &repositoryMock{
				existsFunc: func(ctx context.Context, username, email string) (bool, error) {
					return false, errors.New("some error")
				},
			},
			expectedUser:  nil,
			expectedError: fmt.Errorf("could not check if user exists: some error"),
		},
		{
			name:      "user is created",
			givenUser: givenUser,
			givenRepoMock: &repositoryMock{
				existsFunc: func(ctx context.Context, username, email string) (bool, error) {
					return false, nil
				},
				insertFunc: func(ctx context.Context, user *repository.User) (*repository.User, error) {
					assert.NotEmpty(t, user.ID)
					assert.NotEmpty(t, user.Hash)
					assert.NotEmpty(t, user.CreatedAt)
					assert.NotEmpty(t, user.UpdatedAt)

					return &repository.User{
						ID:        "123",
						Fullname:  givenUser.Fullname,
						Username:  givenUser.Username,
						Birthdate: givenUser.Birthdate,
						Email:     givenUser.Email,
						Hash:      givenUser.Password,
						Role:      string(RoleUser),
						CreatedAt: time.Time{}.AddDate(2000, 1, 1),
						UpdatedAt: time.Time{}.AddDate(2000, 2, 2),
					}, nil
				},
			},
			expectedUser: &User{
				ID:        "123",
				Fullname:  givenUser.Fullname,
				Username:  givenUser.Username,
				Birthdate: givenUser.Birthdate,
				Email:     givenUser.Email,
				Role:      RoleUser,
				CreatedAt: time.Time{}.AddDate(2000, 1, 1),
				UpdatedAt: time.Time{}.AddDate(2000, 2, 2),
			},
			expectedError: nil,
		},
		{
			name:      "insert user error",
			givenUser: givenUser,
			givenRepoMock: &repositoryMock{
				existsFunc: func(ctx context.Context, username, email string) (bool, error) {
					return false, nil
				},
				insertFunc: func(ctx context.Context, user *repository.User) (*repository.User, error) {
					return nil, errors.New("some error")
				},
			},
			expectedUser:  nil,
			expectedError: fmt.Errorf("could not insert user: some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := Service{
				repo: tc.givenRepoMock,
			}

			user, err := svc.Create(context.Background(), tc.givenUser)
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedUser, user)
		})
	}
}

func TestFetchByID(t *testing.T) {
	t.Parallel()

	givenID := uuid.New().String()

	testCases := []struct {
		name          string
		givenID       string
		givenRepoMock *repositoryMock
		expectedUser  *User
		expectedError error
	}{
		{
			name:    "user not found",
			givenID: givenID,
			givenRepoMock: &repositoryMock{
				selectByIDFunc: func(ctx context.Context, id string) (*repository.User, error) {
					return nil, nil
				},
			},
			expectedUser:  nil,
			expectedError: errNotFound,
		},
		{
			name:    "select user error",
			givenID: givenID,
			givenRepoMock: &repositoryMock{
				selectByIDFunc: func(ctx context.Context, id string) (*repository.User, error) {
					return nil, errors.New("some error")
				},
			},
			expectedUser:  nil,
			expectedError: fmt.Errorf("could not select user by id: some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := Service{
				repo: tc.givenRepoMock,
			}

			user, err := svc.FetchByID(context.Background(), tc.givenID)
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedUser, user)
		})
	}
}

func TestFetchByID_validation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		givenID       string
		expectedError error
	}{
		{
			name:          "emty id",
			givenID:       "",
			expectedError: errIDEmpty,
		},
		{
			name:          "invalid id",
			givenID:       "%invalid-id%",
			expectedError: errIDInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := Service{}

			_, err := svc.FetchByID(context.Background(), tc.givenID)
			require.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGenerateJWT(t *testing.T) {
	t.Parallel()

	givenUserID := "123"
	givenUserRole := RoleAdmin

	svc := Service{
		jwtSigningKey: "secret",
	}

	actual, err := svc.generateJWT(givenUserID, givenUserRole)
	require.NoError(t, err)

	// Verify the token
	token, err := jwt.Parse(actual, func(token *jwt.Token) (interface{}, error) {
		return []byte(svc.jwtSigningKey), nil
	})
	require.NoError(t, err)

	// Verify the claims
	claims := token.Claims.(jwt.MapClaims)
	assert.Equal(t, givenUserID, claims["id"])
	assert.Equal(t, string(givenUserRole), claims["role"])

	// Verify the expiration
	exp, ok := claims["exp"].(float64)
	require.True(t, ok)

	assert.True(t, exp > float64(time.Now().Unix()))

	// Verify the signature
	alg, ok := token.Method.(*jwt.SigningMethodHMAC)
	require.True(t, ok)

	assert.Equal(t, jwt.SigningMethodHS256.Alg(), alg.Alg())
}

func TestNewUserFromRepository(t *testing.T) {
	t.Parallel()

	given := repository.User{
		ID:        "123",
		Fullname:  "John Doe",
		Username:  "jdoe",
		Birthdate: "2000-01-01",
		Email:     "jdoe@mail.com",
		Hash:      "password",
		Role:      "admin",
		CreatedAt: time.Time{}.AddDate(2000, 1, 1),
		UpdatedAt: time.Time{}.AddDate(2000, 2, 2),
	}

	expected := User{
		ID:        "123",
		Fullname:  "John Doe",
		Username:  "jdoe",
		Birthdate: "2000-01-01",
		Email:     "jdoe@mail.com",
		Role:      RoleAdmin,
		CreatedAt: time.Time{}.AddDate(2000, 1, 1),
		UpdatedAt: time.Time{}.AddDate(2000, 2, 2),
	}

	actual, err := newUserFromRepository(&given)
	require.NoError(t, err)

	assert.Equal(t, expected, *actual)
}
