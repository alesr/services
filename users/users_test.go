package users

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/alesr/stdservices/internal/users/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestNew(t *testing.T) {
	t.Parallel()

	givenLogger := zap.NewNop()
	givenAppName := "test-app"
	givenJWTSigningKey := "secret"
	givenEmailVerificationSecret := "secret"

	givenEmailVerificationEndpoint, err := url.Parse("http://localhost:8080/verify-email")
	require.NoError(t, err)

	givenEmailer := &emailerMock{}
	givenRepo := &repositoryMock{}

	actual := New(
		givenLogger,
		givenAppName,
		givenJWTSigningKey,
		givenEmailVerificationSecret,
		*givenEmailVerificationEndpoint,
		givenEmailer,
		givenRepo,
	)

	require.NotNil(t, actual)
	require.Equal(t, givenLogger, actual.logger)
	assert.Equal(t, givenAppName, actual.appName)
	assert.Equal(t, givenJWTSigningKey, actual.jwtSigningKey)
	assert.Equal(t, givenEmailVerificationSecret, actual.emailVerificationSecret)
	assert.Equal(t, givenEmailVerificationEndpoint.String(), actual.emailVerificationEndpoint.String())
	assert.Equal(t, givenEmailer, actual.emailer)
	assert.Equal(t, givenRepo, actual.repo)
}

func TestCreate_validation(t *testing.T) {
	t.Parallel()

	given := CreateUserInput{
		Fullname: "%invalid-name%",
	}

	svc := DefaultService{}

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

	givenUserWithAdminRole := givenUser
	givenUserWithAdminRole.Role = RoleAdmin

	testCases := []struct {
		name             string
		givenUser        CreateUserInput
		givenEmailerMock *emailerMock
		givenRepoMock    *repositoryMock
		expectedUser     *User
		expectedError    error
	}{
		{
			name:      "user aleready exists",
			givenUser: givenUser,
			givenRepoMock: &repositoryMock{
				insertFunc: func(ctx context.Context, user *repository.User) (*repository.User, error) {
					return nil, repository.ErrDuplicateRecord
				},
			},
			expectedUser:  nil,
			expectedError: errAlreadyExists,
		},
		{
			name:      "user is created",
			givenUser: givenUser,
			givenEmailerMock: &emailerMock{
				sendFunc: func(to, subject, body string) error {
					return nil
				},
			},
			givenRepoMock: &repositoryMock{
				insertEmailVerificationFunc: func(ctx context.Context, in repository.EmailVerification) error {
					assert.NotEmpty(t, in.UserID)
					assert.NotEmpty(t, in.Token)
					assert.NotEmpty(t, in.CreatedAt)
					assert.NotEmpty(t, in.ExpiresAt)
					return nil
				},
				insertFunc: func(ctx context.Context, user *repository.User) (*repository.User, error) {
					assert.NotEmpty(t, user.ID)
					assert.NotEmpty(t, user.PasswordHash)
					assert.NotEmpty(t, user.CreatedAt)
					assert.NotEmpty(t, user.UpdatedAt)

					return &repository.User{
						ID:           "123",
						Fullname:     givenUser.Fullname,
						Username:     givenUser.Username,
						Birthdate:    givenUser.Birthdate,
						Email:        givenUser.Email,
						PasswordHash: givenUser.Password,
						Role:         string(RoleUser),
						CreatedAt:    time.Time{}.AddDate(2000, 1, 1),
						UpdatedAt:    time.Time{}.AddDate(2000, 2, 2),
					}, nil
				},
			},
			expectedUser: &User{
				ID:            "123",
				Fullname:      givenUser.Fullname,
				Username:      givenUser.Username,
				Birthdate:     givenUser.Birthdate,
				Email:         givenUser.Email,
				EmailVerified: false,
				Role:          RoleUser,
				CreatedAt:     time.Time{}.AddDate(2000, 1, 1),
				UpdatedAt:     time.Time{}.AddDate(2000, 2, 2),
			},
			expectedError: nil,
		},
		{
			name:      "insert user error",
			givenUser: givenUser,
			givenRepoMock: &repositoryMock{
				insertFunc: func(ctx context.Context, user *repository.User) (*repository.User, error) {
					return nil, errors.New("some error")
				},
			},
			expectedUser:  nil,
			expectedError: fmt.Errorf("could not insert user: some error"),
		},
		{
			name:          "cannot create user with admin role",
			givenUser:     givenUserWithAdminRole,
			givenRepoMock: &repositoryMock{},
			expectedUser:  nil,
			expectedError: errCannotCreateAdminUser,
		},
		{
			name:      "send email verification error still creates an user",
			givenUser: givenUser,
			givenEmailerMock: &emailerMock{
				sendFunc: func(to, subject, body string) error {
					return errors.New("some error")
				},
			},
			givenRepoMock: &repositoryMock{
				insertEmailVerificationFunc: func(ctx context.Context, in repository.EmailVerification) error {
					assert.NotEmpty(t, in.UserID)
					assert.NotEmpty(t, in.Token)
					assert.NotEmpty(t, in.CreatedAt)
					assert.NotEmpty(t, in.ExpiresAt)
					return nil
				},
				insertFunc: func(ctx context.Context, user *repository.User) (*repository.User, error) {
					assert.NotEmpty(t, user.ID)
					assert.NotEmpty(t, user.PasswordHash)
					assert.NotEmpty(t, user.CreatedAt)
					assert.NotEmpty(t, user.UpdatedAt)

					return &repository.User{
						ID:           "123",
						Fullname:     givenUser.Fullname,
						Username:     givenUser.Username,
						Birthdate:    givenUser.Birthdate,
						Email:        givenUser.Email,
						PasswordHash: givenUser.Password,
						Role:         string(RoleUser),
						CreatedAt:    time.Time{}.AddDate(2000, 1, 1),
						UpdatedAt:    time.Time{}.AddDate(2000, 2, 2),
					}, nil
				},
			},
			expectedUser: &User{
				ID:            "123",
				Fullname:      givenUser.Fullname,
				Username:      givenUser.Username,
				Birthdate:     givenUser.Birthdate,
				Email:         givenUser.Email,
				EmailVerified: false,
				Role:          RoleUser,
				CreatedAt:     time.Time{}.AddDate(2000, 1, 1),
				UpdatedAt:     time.Time{}.AddDate(2000, 2, 2),
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := DefaultService{
				logger:  zap.NewNop(),
				emailer: tc.givenEmailerMock,
				repo:    tc.givenRepoMock,
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
			svc := DefaultService{
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
		expectedError bool
	}{
		{
			name:          "emty id",
			givenID:       "",
			expectedError: true,
		},
		{
			name:          "invalid id",
			givenID:       "%invalid-id%",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := DefaultService{}

			_, err := svc.FetchByID(context.Background(), tc.givenID)
			require.Equal(t, tc.expectedError, err != nil)
		})
	}
}

func TestDelete_validation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		givenID       string
		expectedError bool
	}{
		{
			name:          "emty id",
			givenID:       "",
			expectedError: true,
		},
		{
			name:          "invalid id",
			givenID:       "%invalid-id%",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := DefaultService{}

			err := svc.Delete(context.Background(), tc.givenID)
			require.Equal(t, tc.expectedError, err != nil)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		givenRepoMock *repositoryMock
		expectError   bool
	}{
		{
			name: "delete user error",
			givenRepoMock: &repositoryMock{
				deleteByIDFunc: func(ctx context.Context, id string) error {
					return errors.New("some error")
				},
			},
			expectError: true,
		},
		{
			name: "delete user success",
			givenRepoMock: &repositoryMock{
				deleteByIDFunc: func(ctx context.Context, id string) error {
					return nil
				},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := DefaultService{
				repo: tc.givenRepoMock,
			}

			err := svc.Delete(context.Background(), uuid.New().String())
			assert.Equal(t, tc.expectError, err != nil)
		})
	}
}

func TestGenerateToken_validation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		givenEmail    string
		givenPassword string
		expectedError bool
	}{
		{
			name:          "empty email",
			givenEmail:    "",
			givenPassword: "password%&123",
			expectedError: true,
		},
		{
			name:          "invalid email format",
			givenEmail:    "invalid-email",
			givenPassword: "password%&123",
			expectedError: true,
		},
		{
			name:          "empty password",
			givenEmail:    "joedoe@mail.com",
			givenPassword: "",
			expectedError: true,
		},
		{
			name:          "empty password format",
			givenEmail:    "joedoe@mail.com",
			givenPassword: "123",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := DefaultService{}

			_, err := svc.GenerateToken(context.Background(), tc.givenEmail, tc.givenPassword)
			require.Equal(t, tc.expectedError, err != nil)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	t.Parallel()

	password := "password%&123"

	givenHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		givenPassword string
		givenRepoMock *repositoryMock
		expectedToken bool
		expectedError bool
	}{
		{
			name:          "user not found",
			givenPassword: password,
			givenRepoMock: &repositoryMock{
				selectByEmailFunc: func(ctx context.Context, email string) (*repository.User, error) {
					return nil, nil
				},
			},
			expectedToken: false,
			expectedError: true,
		},
		{
			name:          "select user error",
			givenPassword: password,
			givenRepoMock: &repositoryMock{
				selectByEmailFunc: func(ctx context.Context, email string) (*repository.User, error) {
					return nil, errors.New("some error")
				},
			},
			expectedToken: false,
			expectedError: true,
		},
		{
			name:          "password match",
			givenPassword: password,
			givenRepoMock: &repositoryMock{
				selectByEmailFunc: func(ctx context.Context, email string) (*repository.User, error) {
					return &repository.User{
						ID:           uuid.New().String(),
						Role:         string(RoleUser),
						Email:        email,
						PasswordHash: string(givenHash),
					}, nil
				},
			},
			expectedToken: true,
			expectedError: false,
		},
		{
			name:          "password not match",
			givenPassword: "somepassword&#%123",
			givenRepoMock: &repositoryMock{
				selectByEmailFunc: func(ctx context.Context, email string) (*repository.User, error) {
					return &repository.User{
						ID:           uuid.New().String(),
						Role:         string(RoleUser),
						Email:        email,
						PasswordHash: string(givenHash),
					}, nil
				},
			},
			expectedToken: false,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := DefaultService{
				repo: tc.givenRepoMock,
			}

			token, err := svc.GenerateToken(context.Background(), "joedoe@mail.com", tc.givenPassword)
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, tc.expectedToken, token != "")
		})
	}
}

func TestNewUserFromRepository(t *testing.T) {
	t.Parallel()

	givenUser := repository.User{
		ID:            "123",
		Fullname:      "John Doe",
		Username:      "jdoe",
		Birthdate:     "2000-01-01",
		Email:         "jdoe@mail.com",
		EmailVerified: true,
		PasswordHash:  "password",
		CreatedAt:     time.Time{}.AddDate(2000, 1, 1),
		UpdatedAt:     time.Time{}.AddDate(2000, 2, 2),
	}

	t.Run("valid role admin", func(t *testing.T) {
		given := givenUser
		given.Role = string(RoleAdmin)

		expected := User{
			ID:            "123",
			Fullname:      "John Doe",
			Username:      "jdoe",
			Birthdate:     "2000-01-01",
			Email:         "jdoe@mail.com",
			EmailVerified: true,
			Role:          RoleAdmin,
			CreatedAt:     time.Time{}.AddDate(2000, 1, 1),
			UpdatedAt:     time.Time{}.AddDate(2000, 2, 2),
		}

		actual, err := newUserFromRepository(&given)
		require.NoError(t, err)

		assert.Equal(t, expected, *actual)
	})

	t.Run("valid role user", func(t *testing.T) {
		given := givenUser
		given.Role = string(RoleUser)

		expected := User{
			ID:            "123",
			Fullname:      "John Doe",
			Username:      "jdoe",
			Birthdate:     "2000-01-01",
			Email:         "jdoe@mail.com",
			EmailVerified: true,
			Role:          RoleUser,
			CreatedAt:     time.Time{}.AddDate(2000, 1, 1),
			UpdatedAt:     time.Time{}.AddDate(2000, 2, 2),
		}

		actual, err := newUserFromRepository(&given)
		require.NoError(t, err)

		assert.Equal(t, expected, *actual)
	})

	t.Run("invalid role", func(t *testing.T) {
		given := givenUser
		given.Role = "invalid"

		_, err := newUserFromRepository(&given)
		assert.Error(t, err)
	})
}
