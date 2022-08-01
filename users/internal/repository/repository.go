package repository

import (
	"context"
	"time"
)

type (
	Repository interface {
		Exists(ctx context.Context, username, email string) (bool, error)
		Insert(ctx context.Context, user *User) (*User, error)
		SelectByID(ctx context.Context, id string) (*User, error)
		SelectByEmail(ctx context.Context, email string) (*User, error)
	}

	// User represents a user in the database table
	User struct {
		ID        string
		Firstname string
		Lastname  string
		Username  string
		Birthdate string
		Email     string
		Hash      string
		Role      string
		CreatedAt time.Time
	}
)
