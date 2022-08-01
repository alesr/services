package validate

import "errors"

var (
	// List error messages

	errBirthdateFormat   = errors.New("birthdate must be in the format YYYY-MM-DD")
	errBirthdateRequired = errors.New("birthdate is required")
	errEmailFormat       = errors.New("email is invalid")
	errEmailRequired     = errors.New("email is required")
	errFullnameFormat    = errors.New("fullname must only contain letters and spaces")
	errFullnameLength    = errors.New("fullname must be between 3 and 64 characters")
	errFullnameRequired  = errors.New("fullname is required")
	errPasswordFormat    = errors.New("password must contain at least one number, one letter and one special character")
	errPasswordLength    = errors.New("password must be between 8 and 64 characters")
	errPasswordRequired  = errors.New("password is required")
)
