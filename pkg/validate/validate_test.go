package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullname(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		given    string
		expected error
	}{
		{
			name:     "valid",
			given:    "John Doe",
			expected: nil,
		},
		{
			name:     "empty",
			given:    "",
			expected: errFullnameRequired,
		},
		{
			name:     "too short",
			given:    "a",
			expected: errFullnameLength,
		},
		{
			name:     "too long",
			given:    "mckrbdwenwfrbvkgivwqivchjvvijvuycprqdnddjqdnnfwiczwhrfxznnzxpnmjl",
			expected: errFullnameLength,
		},
		{
			name:     "invalid characters",
			given:    "a*",
			expected: errFullnameFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Fullname(tc.given)
			assert.Equal(t, actual, tc.expected)
		})
	}
}

func TestBirthdate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		given    string
		expected error
	}{
		{
			name:     "valid",
			given:    "1990-01-01",
			expected: nil,
		},
		{
			name:     "empty",
			given:    "",
			expected: errBirthdateRequired,
		},
		{
			name:     "invalid format",
			given:    "2019/01/01",
			expected: errBirthdateFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Birthdate(tc.given)
			assert.Equal(t, actual, tc.expected)
		})
	}
}

func TestEmail(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		given    string
		expected error
	}{
		{
			name:     "valid",
			given:    "jdoe@mail.com",
			expected: nil,
		},
		{
			name:     "empty",
			given:    "",
			expected: errEmailRequired,
		},
		{
			name:     "invalid format",
			given:    "abc",
			expected: errEmailFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Email(tc.given)
			assert.Equal(t, actual, tc.expected)
		})
	}
}

func TestPassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		given    string
		expected error
	}{
		{
			name:     "valid",
			given:    "abcd01!@",
			expected: nil,
		},
		{
			name:     "empty",
			given:    "",
			expected: errPasswordRequired,
		},
		{
			name:     "too short",
			given:    "a",
			expected: errPasswordLength,
		},
		{
			name:     "too long",
			given:    "mckrbdwenwfrbvkgivwqivchjvvijvuycprqdnddjqdnnfwiczwhrfxznnzxpnmjl",
			expected: errPasswordLength,
		},
		{
			name:     "only letters",
			given:    "abcdefghijklmnopqrstuvwxyz",
			expected: errPasswordFormat,
		},
		{
			name:     "only numbers",
			given:    "0123456789",
			expected: errPasswordFormat,
		},
		{
			name:     "only special characters",
			given:    "!@#$%^&*()_+-=",
			expected: errPasswordFormat,
		},
		{
			name:     "only letters and numbers",
			given:    "abcdefghijklmnopqrstuvwxyz0123456789",
			expected: errPasswordFormat,
		},
		{
			name:     "only letters and special characters",
			given:    "abcdefghijklmnopqrstuvwxyz!@#$%^&*()_+-=",
			expected: errPasswordFormat,
		},
		{
			name:     "only numbers and special characters",
			given:    "0123456789!@#$%^&*()_+-=",
			expected: errPasswordFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Password(tc.given)
			assert.Equal(t, actual, tc.expected)
		})
	}
}
