package repository

import (
	"time"
)

// User represents a user in the database table
type User struct {
	ID        string
	Fullname  string
	Username  string
	Birthdate string
	Email     string
	Hash      string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
