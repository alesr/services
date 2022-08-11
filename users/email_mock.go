package users

import "errors"

var _ emailer = (*emailerMock)(nil)

type emailerMock struct {
	sendFunc func(to string, body []byte) error
}

func (m *emailerMock) Send(to string, body []byte) error {
	if m.sendFunc == nil {
		return errors.New("emailerMock.sendFunc is nil")
	}
	return m.sendFunc(to, body)
}
