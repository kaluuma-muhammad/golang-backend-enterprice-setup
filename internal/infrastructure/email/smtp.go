package email

import (
	"context"
	"fmt"

	mail "github.com/wneessen/go-mail"
)

type SMTPProvider struct {
	client      *mail.Client
	fromAddress string
	fromName    string
}

func NewSMTP(host string, port int, username, password, fromAddress, fromName string, encryption string) (*SMTPProvider, error) {
	options := []mail.Option{
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	}

	switch encryption {
	case "ssl":
		options = append(options, mail.WithSSL())

	case "tls":
		options = append(options, mail.WithTLSPolicy(mail.TLSMandatory))

	default:
		options = append(options, mail.WithTLSPolicy(mail.TLSOpportunistic))
	}

	client, err := mail.NewClient(
		host,
		options...,
	)
	if err != nil {
		return nil, err
	}

	return &SMTPProvider{
		client:      client,
		fromAddress: fromAddress,
		fromName:    fromName,
	}, nil
}

func (s *SMTPProvider) Send(ctx context.Context, message Message) error {
	msg := mail.NewMsg()

	if err := msg.FromFormat(s.fromName, s.fromAddress); err != nil {
		return err
	}

	if err := msg.To(message.To); err != nil {
		return err
	}

	msg.Subject(message.Subject)

	if message.HTML != "" {
		msg.AddAlternativeString(mail.TypeTextHTML, message.HTML)
	}

	if err := s.client.DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
