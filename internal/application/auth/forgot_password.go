package auth

import (
	"context"
	"time"

	backendSession "github.com/go-api/internal/domain/session"
	domainToken "github.com/go-api/internal/domain/token"
)

func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return ErrUserNotFound
	}

	tokenType := domainToken.PasswordReset
	code, err := s.tokenService.CreateCodeToken(ctx, user.ID, tokenType)
	if err != nil {
		return err
	}

	go s.emailService.SendPasswordResetEmail(context.Background(), user, code)

	return nil
}

func (s *Service) ResendVerification(ctx context.Context, req ResendVerificationRequest) error {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return ErrUserNotFound
	}

	tokenType := domainToken.PasswordReset
	code, err := s.tokenService.CreateCodeToken(ctx, user.ID, tokenType)
	if err != nil {
		return err
	}

	go s.emailService.SendVerificationEmail(context.Background(), user, code)

	return nil
}

func (s *Service) VerifyResetCode(ctx context.Context, req VerifyResetCodeRequest) (*VerifyResetCodeResponse, error) {
	var response *VerifyResetCodeResponse

	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		user, err := s.users.FindByEmail(txCtx, req.Email)
		if err != nil {
			return ErrInvalidVerificationCode
		}

		hash := s.tokenService.Hash(req.Code)
		token, err := s.tokenService.FindByTokenAndTypeAndUser(
			txCtx,
			hash,
			domainToken.PasswordReset,
			user.ID,
		)

		if err != nil {
			return ErrInvalidVerificationCode
		}

		if token.UserID != user.ID {
			return ErrInvalidVerificationCode
		}

		if token.UsedAt != nil {
			return ErrInvalidVerificationCode
		}

		if token.ExpiresAt.Before(time.Now()) {
			_ = s.tokenService.DeleteByUserAndType(txCtx, user.ID, domainToken.PasswordReset)

			return ErrVerificationCodeExpired
		}

		err = s.tokenService.Consume(txCtx, token.ID)
		if err != nil {
			return err
		}

		resetToken, err := s.tokenService.CreateResetToken(txCtx, user.ID)
		if err != nil {
			return err
		}

		response = &VerifyResetCodeResponse{
			ResetToken: resetToken,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		hash := s.tokenService.Hash(req.ResetToken)

		token, err := s.tokenService.FindByTokenAndType(
			txCtx,
			hash,
			domainToken.PasswordResetGrant,
		)
		if err != nil {
			return ErrInvalidResetToken
		}

		if token.UsedAt != nil {
			return ErrResetTokenAlreadyUsed
		}

		if token.ExpiresAt.Before(time.Now()) {
			_ = s.tokenService.DeleteByUserAndType(txCtx, token.UserID, domainToken.PasswordResetGrant)

			return ErrResetTokenExpired
		}

		user, err := s.users.FindByID(txCtx, token.UserID)
		if err != nil {
			return ErrInvalidResetToken
		}

		password, err := s.password.Hash(req.NewPassword)
		if err != nil {
			return err
		}

		if err := s.users.UpdatePassword(txCtx, user.ID, password); err != nil {
			return err
		}

		if err := s.tokenService.Consume(txCtx, token.ID); err != nil {
			return err
		}

		if err := s.sessions.RevokeAllSessions(
			txCtx,
			user.ID,
			backendSession.RevokedByPasswordChange,
		); err != nil {
			return err
		}
		return nil
	})

	return err
}
