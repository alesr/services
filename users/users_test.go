package users

import (
	"testing"
	"time"

	"github.com/alesr/services/users/internal/repository"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
