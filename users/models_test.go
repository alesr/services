package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserInput_validate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		given         CreateUserInput
		expectedError bool
	}{
		{
			name: "valid",
			given: CreateUserInput{
				Fullname:        "John Doe",
				Username:        "johndoe",
				Birthdate:       "1990-01-01",
				Email:           "joedoe@mail.com",
				Password:        "1234%6abc",
				ConfirmPassword: "1234%6abc",
				Role:            string(RoleAdmin),
			},
			expectedError: false,
		},
		{
			name: "missing fullname",
			given: CreateUserInput{
				Fullname:        "",
				Username:        "johndoe",
				Birthdate:       "1990-01-01",
				Email:           "joedoe@mail.com",
				Password:        "1234%6abc",
				ConfirmPassword: "1234%6abc",
				Role:            string(RoleAdmin),
			},
			expectedError: true,
		},
		{
			name: "missing username",
			given: CreateUserInput{
				Fullname:        "John Doe",
				Username:        "",
				Birthdate:       "1990-01-01",
				Email:           "joedoe@mail.com",
				Password:        "1234%6abc",
				ConfirmPassword: "1234%6abc",
				Role:            string(RoleAdmin),
			},
			expectedError: true,
		},
		{
			name: "missing birthdate",
			given: CreateUserInput{
				Fullname:        "John Doe",
				Username:        "johndoe",
				Birthdate:       "",
				Email:           "joedoe@mail.com",
				Password:        "1234%6abc",
				ConfirmPassword: "1234%6abc",
				Role:            string(RoleAdmin),
			},
			expectedError: true,
		},
		{
			name: "missing email",
			given: CreateUserInput{
				Fullname:        "John Doe",
				Username:        "johndoe",
				Birthdate:       "1990-01-01",
				Email:           "",
				Password:        "1234%6abc",
				ConfirmPassword: "1234%6abc",
				Role:            string(RoleAdmin),
			},
			expectedError: true,
		},
		{
			name: "missing password",
			given: CreateUserInput{
				Fullname:        "John Doe",
				Username:        "johndoe",
				Birthdate:       "1990-01-01",
				Email:           "joedoe@mail.com",
				Password:        "",
				ConfirmPassword: "1234%6abc",
				Role:            string(RoleAdmin),
			},
			expectedError: true,
		},
		{
			name: "missing confirm password",
			given: CreateUserInput{
				Fullname:        "John Doe",
				Username:        "johndoe",
				Birthdate:       "1990-01-01",
				Email:           "joedoe@mail.com",
				Password:        "1234%6abc",
				ConfirmPassword: "",
				Role:            string(RoleAdmin),
			},
			expectedError: true,
		},
		{
			name: "password mismatch",
			given: CreateUserInput{
				Fullname:        "John Doe",
				Username:        "johndoe",
				Birthdate:       "1990-01-01",
				Email:           "joedoe@mail.com",
				Password:        "1234%6abc",
				ConfirmPassword: "1234%6zzzz",
				Role:            string(RoleAdmin),
			},
			expectedError: true,
		},
		{
			name: "invalid role",
			given: CreateUserInput{
				Fullname:        "John Doe",
				Username:        "johndoe",
				Birthdate:       "1990-01-01",
				Email:           "joedoe@mail.com",
				Password:        "1234%6abc",
				ConfirmPassword: "1234%6abc",
				Role:            "invalid",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.given.validate()

			if tc.expectedError {
				assert.Error(t, actual)
			} else {
				assert.NoError(t, actual)
			}
		})
	}
}
