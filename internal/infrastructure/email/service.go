package email

import (
	"context"
	"errors"
)

type Service struct {
	provider Provider
}

func New(provider Provider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) Send(ctx context.Context, message Message) error {
	if message.To == "" {
		return errors.New("recipient email is required")
	}

	if message.Subject == "" {
		return errors.New("email subject is required")
	}

	if message.HTML == "" && message.Text == "" {
		return errors.New("email body must contain either HTML or plain text")
	}

	return s.provider.Send(ctx, message)
}
