package interactors

import "io"

type Emailer interface {
	EmailFile(r io.Reader, name, to string) error
}

type SendFile struct {
	e Emailer
}

func NewSendFile(e Emailer) *SendFile {
	return &SendFile{
		e: e,
	}
}

func (s *SendFile) Execute(r io.Reader, name, toEmail string) error {
	return s.e.EmailFile(r, name, toEmail)
}
