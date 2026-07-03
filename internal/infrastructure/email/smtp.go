package email

import "context"

type SMTPProvider struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewSMTP(
	host,
	port,
	user,
	pass,
	from string,
) *SMTPProvider {
	return &SMTPProvider{
		host: host,
		port: port,
		user: user,
		pass: pass,
		from: from,
	}
}

func (s *SMTPProvider) Send(
	ctx context.Context,
	message Message,
) error {

	return nil
}
