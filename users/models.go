package users

import (
	"fmt"
	"time"
	"unicode"
)

const birthdateFormat string = "2006-01-02"

type (
	// User represents a user domain model
	User struct {
		ID        string
		Firstname string
		Lastname  string
		Username  string
		Birthdate string
		Email     string
		Role      role
		CreatedAt time.Time
	}

	// CreateUserInput represents the input data for creating a user
	CreateUserInput struct {
		Firstname       string
		Lastname        string
		Username        string
		Birthdate       string
		Email           string
		Password        string
		ConfirmPassword string
		Role            role
	}
)

func (in *CreateUserInput) validate() error {
	if err := validateName(in.Firstname, 3, 20); err != nil {
		return fmt.Errorf("could not validate first name: %w", err)
	}

	if err := validateName(in.Lastname, 3, 20); err != nil {
		return fmt.Errorf("could not validate last name: %w", err)
	}

	if err := validateName(in.Username, 3, 20); err != nil {
		return fmt.Errorf("could not validate username: %w", err)
	}

	if in.Birthdate == "" {
		return ErrBirthdateRequired
	}

	if _, err := time.Parse(birthdateFormat, in.Birthdate); err != nil {
		return ErrInvalidBirthdate
	}

	if in.Email == "" {
		return ErrEmailRequired
	}

	if len(in.Email) > 255 {
		return ErrEmailTooLong
	}

	if in.Password == "" {
		return ErrPasswordRequired
	}

	if len(in.Password) < 6 {
		return ErrPasswordTooShort
	}

	if len(in.Password) > 128 {
		return ErrPasswordTooLong
	}

	if in.Password != in.ConfirmPassword {
		return ErrPasswordMismatch
	}

	if in.Password != in.ConfirmPassword {
		return ErrPasswordMismatch
	}
	return nil
}

type AuthUserInput struct {
	Email    string
	Password string
}

func (in *AuthUserInput) validate() error {
	if in.Email == "" {
		return ErrEmailRequired
	}

	if len(in.Email) > 255 {
		return ErrEmailTooLong
	}

	if in.Password == "" {
		return ErrPasswordRequired
	}

	if len(in.Password) < 6 {
		return ErrPasswordTooShort
	}

	if len(in.Password) > 128 {
		return ErrPasswordTooLong
	}
	return nil
}

func validateName(name string, minLen, maxLen int) error {
	if name == "" {
		return ErrNameRequired
	}

	if len(name) < minLen {
		return ErrNameTooShort
	}

	if len(name) > maxLen {
		return ErrNameTooLong
	}

	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return ErrNameInvalid
		}
	}
	return nil
}
