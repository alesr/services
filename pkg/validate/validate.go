package validate

import (
	"net/mail"
	"time"
	"unicode"

	"github.com/google/uuid"
)

const (
	minFullnameLen = 3
	maxFullnameLen = 64

	birthdateFormat string = "2006-01-02"
)

func Fullname(name string) error {
	if name == "" {
		return errFullnameRequired
	}

	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return errFullnameFormat
		}
	}

	if len(name) < minFullnameLen || len(name) > maxFullnameLen {
		return errFullnameLength
	}
	return nil
}

func Birthdate(bDate string) error {
	if bDate == "" {
		return errBirthdateRequired
	}

	// check birthdate has format 2006-01-02
	if _, err := time.Parse(birthdateFormat, bDate); err != nil {
		return errBirthdateFormat
	}
	return nil
}

func Email(email string) error {
	if email == "" {
		return errEmailRequired
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return errEmailFormat
	}
	return nil
}

func Password(password string) error {
	if password == "" {
		return errPasswordRequired
	}

	// check password has length between 8 and 64
	if len(password) < 8 || len(password) > 64 {
		return errPasswordLength
	}

	var hasNumber, hasLetter, hasSpecial bool

	for _, char := range password {
		if unicode.IsNumber(char) {
			hasNumber = true
		}
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			hasSpecial = true
		}
	}

	if !hasNumber || !hasLetter || !hasSpecial {
		return errPasswordFormat
	}
	return nil
}

func ID(id string) error {
	if id == "" {
		return errIDRequired
	}

	if _, err := uuid.Parse(id); err != nil {
		return errIDFormat
	}
	return nil
}
