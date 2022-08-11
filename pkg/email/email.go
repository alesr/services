package email

import (
	"fmt"
	"net/smtp"
)

type Config struct {
	Sender   string
	Identity string
	Username string
	Password string
	Host     string
	Port     string
}

type mailer struct {
	sender string
	host   string
	port   string
	auth   smtp.Auth
}

func New(cfg Config) *mailer {
	return &mailer{
		sender: cfg.Sender,
		host:   cfg.Host,
		port:   cfg.Port,
		auth:   smtp.PlainAuth(cfg.Identity, cfg.Username, cfg.Password, cfg.Host),
	}
}

func (m *mailer) Send(to string, body []byte) error {
	if err := smtp.SendMail(
		m.host+":"+m.port, m.auth, m.sender, []string{to}, body,
	); err != nil {
		return fmt.Errorf("could not send mail: %s", err)
	}
	return nil
}
