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

	errAlreadyExists         = newE("user already exists")
	errCannotCreateAdminUser = newE("cannot create admin user")
	errNotFound              = newE("user not found")
	errPasswordInvalid       = newE("user password is invalid")
	errPasswordMismatch      = newE("user password mismatch")
	errRoleInvalid           = newE("user role is invalid")
	errTokenEmpty            = newE("user token is empty")
	errTokenExpired          = newE("user token is expired")
	errTokenInvalid          = newE("user token is invalid")
)
