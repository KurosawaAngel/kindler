package email

import (
	"io"
	"mime"

	"gopkg.in/gomail.v2"
)

type Config struct {
	Host     string `env:"SMTP_HOST" env-required:"true"`
	Port     int    `env:"SMTP_PORT" env-required:"true"`
	Username string `env:"SMTP_USERNAME" env-required:"true"`
	Password string `env:"SMTP_PASSWORD" env-required:"true"`
}

type Emailer struct {
	g *gomail.Dialer
}

func New(g *gomail.Dialer) *Emailer {
	return &Emailer{
		g: g,
	}
}

func (e *Emailer) EmailFile(data io.Reader, name string, to string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", e.g.Username)
	m.SetHeader("To", to)
	m.SetHeader("Content-Type", "multipart/mixed")
	m.SetHeader("Subject", mime.QEncoding.Encode("UTF-8", name))
	m.SetBody("text/plain; charset=UTF-8", name)

	m.Attach(name, gomail.Rename(mime.QEncoding.Encode("UTF-8", name)),
		gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, data)
			return err
		}))

	return e.g.DialAndSend(m)
}
