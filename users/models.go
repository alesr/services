package users

import (
	"time"

	"github.com/alesr/services/pkg/validate"
)

const (
	// Enumerate available roles

	RoleAdmin role = iota + 1
	RoleUser
)

type AuthUserInput struct {
	Email    string
	Password string
}

type role uint8

func (r role) validate() error {
	if r != RoleAdmin && r != RoleUser {
		return errRoleInvalid
	}
	return nil
}

// User represents a user domain model
type User struct {
	ID        string
	Fullname  string
	Username  string
	Birthdate string
	Email     string
	Role      role
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUserInput represents the input data for creating a user
type CreateUserInput struct {
	Fullname        string
	Username        string
	Birthdate       string
	Email           string
	Password        string
	ConfirmPassword string
	Role            role
}

func (in *CreateUserInput) validate() error {
	if err := validate.Fullname(in.Fullname); err != nil {
		return newE(err.Error())
	}

	if err := validate.Fullname(in.Username); err != nil {
		return newE(err.Error())
	}

	if err := validate.Birthdate(in.Birthdate); err != nil {
		return newE(err.Error())
	}

	if err := validate.Email(in.Email); err != nil {
		return newE(err.Error())
	}

	if err := validate.Password(in.Password); err != nil {
		return newE(err.Error())
	}

	if in.Password != in.ConfirmPassword {
		return errPasswordMissmatch
	}

	if err := in.Role.validate(); err != nil {
		return newE(err.Error())
	}
	return nil
}
