package repository

import (
	"errors"
	"time"
)

var ErrDuplicateRecord error = errors.New("duplicate record")

// User represents a user in the database table
type User struct {
	ID           string
	Fullname     string
	Username     string
	Birthdate    string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
