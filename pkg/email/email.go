package email

import (
	"fmt"
	"net"
	"net/smtp"
)

type email struct {
	auth smtp.Auth
	addr string
}

func New(identity, username, password, host, port string) *email {
	return &email{
		auth: smtp.PlainAuth(identity, username, password, host),
		addr: net.JoinHostPort(host, port),
	}
}

func (e *email) Send(from, to string, body []byte) error {
	if err := smtp.SendMail(e.addr, e.auth, from, []string{to}, body); err != nil {
		return fmt.Errorf("could not send mail: %s", err)
	}
	return nil
}
