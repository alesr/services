package users

const (
	// Enumerate service error codes

	ErrCodeBirthdateInvalid = iota + 1
	ErrCodeBirthdateRequired
	ErrCodeEmailRequired
	ErrCodeEmailTooLong
	ErrCodeNameInvalid
	ErrCodeNameRequired
	ErrCodeNameTooLong
	ErrCodeNameTooShort
	ErrCodePasswordInvalid
	ErrCodePasswordMismatch
	ErrCodePasswordRequired
	ErrCodePasswordTooLong
	ErrCodePasswordTooShort
	ErrCodeTokenExpired
	ErrCodeUserAlreadyExists
	ErrCodeUsernameTaken
	ErrCodeUserNotFound
)

var (
	// Enumerate service errors

	ErrBirthdateRequired      Error = Error{ErrCodeBirthdateRequired, "birth date is required"}
	ErrEmailRequired          Error = Error{ErrCodeEmailRequired, "email is required"}
	ErrEmailTooLong           Error = Error{ErrCodeEmailTooLong, "email is too long"}
	ErrInvalidBirthdate       Error = Error{ErrCodeBirthdateInvalid, "invalid birth date"}
	ErrNameInvalid            Error = Error{ErrCodeNameInvalid, "name is invalid"}
	ErrNameRequired           Error = Error{ErrCodeNameRequired, "name is required"}
	ErrNameTooLong            Error = Error{ErrCodeNameTooLong, "name is too long"}
	ErrNameTooShort           Error = Error{ErrCodeNameTooShort, "name is too short"}
	ErrPasswordInvalid        Error = Error{ErrCodePasswordInvalid, "password is invalid"}
	ErrPasswordMismatch       Error = Error{ErrCodePasswordMismatch, "password and confirm password do not match"}
	ErrPasswordRequired       Error = Error{ErrCodePasswordRequired, "password is required"}
	ErrPasswordTooLong        Error = Error{ErrCodePasswordTooLong, "password is too long"}
	ErrPasswordTooShort       Error = Error{ErrCodePasswordTooShort, "password is too short"}
	ErrTokenExpired           Error = Error{ErrCodeTokenExpired, "token expired"}
	ErrUserAlreadyExistsError Error = Error{ErrCodeUserAlreadyExists, "user already exists"}
	ErrUsernameTaken          Error = Error{ErrCodeUsernameTaken, "username is already taken"}
	ErrUserNotFound           Error = Error{ErrCodeUserNotFound, "user not found"}
)

// Error represents a service error with a code and message
type Error struct {
	code    int
	message string
}

// Error returns the error message implementing the error interface
func (e Error) Error() string {
	return e.message
}

// Code returns the error code
func (e *Error) Code() int {
	return e.code
}
