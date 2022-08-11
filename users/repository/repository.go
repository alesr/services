package repository

import (
	"errors"
	"time"
)

var (
	ErrDuplicateRecord error = errors.New("duplicate record")
	ErrRecordNotFound  error = errors.New("record not found")
)

// User represents a user in the database table
type User struct {
	ID            string
	Fullname      string
	Username      string
	Birthdate     string
	Email         string
	PasswordHash  string
	Role          string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type EmailVerification struct {
	Code      string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}
