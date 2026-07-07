package email

import (
	"context"
	"time"

	backendUser "github.com/go-api/internal/domain/user"
	emailData "github.com/go-api/internal/infrastructure/email/data"
)

func (s *Service) SendPasswordResetEmail(ctx context.Context, user *backendUser.User, code string) error {
	data := emailData.CodeEmail{
		Title:   SubjectPasswordReset,
		AppName: s.cfg.App.Name,
		Name:    user.FirstName,
		Code:    code,
		Expiry:  s.cfg.JWT.PasswordResetMinutes,
		Year:    time.Now().Year(),
	}

	html, err := s.renderer.RenderHTML("password_reset.html", data)
	if err != nil {
		return err
	}

	return s.provider.Send(ctx, Message{
		To:      user.Email,
		Subject: SubjectPasswordReset,
		HTML:    html,
	})
}
