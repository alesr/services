package users

import "errors"

var _ emailer = (*emailerMock)(nil)

type emailerMock struct {
	sendFunc func(to, subject, body string) error
}

func (m *emailerMock) Send(to, subject, body string) error {
	if m.sendFunc == nil {
		return errors.New("emailerMock.sendFunc is nil")
	}
	return m.sendFunc(to, subject, body)
}
