package users

type ParsableError interface {
	Error() string
}

type E struct {
	msg string
}

func newE(msg string) E {
	return E{msg: msg}
}

func (e E) Error() string {
	return e.msg
}

var (
	// Enumerate service errors

	errAlreadyExists     = newE("user already exists")
	errEmailEmpty        = newE("user email is empty")
	errIDEmpty           = newE("user id is empty")
	errIDInvalid         = newE("user id is invalid")
	errNotFound          = newE("user not found")
	errPasswordEmpty     = newE("user password is empty")
	errPasswordInvalid   = newE("user password is invalid")
	errPasswordMissmatch = newE("user password missmatch")
	errRoleInvalid       = newE("user role is invalid")
	errTokenEmpty        = newE("user token is empty")
	errTokenExpired      = newE("user token is expired")
	errTokenInvalid      = newE("user token is invalid")
)
