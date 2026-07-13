package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthServiceContract interface {
	Register(
		ctx context.Context,
		req RegisterRequest,
		userAgent string,
		ipAddress string,
	) (*LoginResponse, error)

	Login(
		ctx context.Context,
		req LoginRequest,
		userAgent string,
		ipAddress string,
	) (*LoginResponse, error)

	ActivateAccount(
		ctx context.Context,
		userID uuid.UUID,
		req ActivateAccountRequest,
	) error

	ForgotPassword(
		ctx context.Context,
		req ForgotPasswordRequest,
	) error

	ResendVerification(
		ctx context.Context,
		req ResendVerificationRequest,
	) error

	VerifyResetCode(
		ctx context.Context,
		req VerifyResetCodeRequest,
	) (*VerifyResetCodeResponse, error)

	ResetPassword(
		ctx context.Context,
		req ResetPasswordRequest,
	) error

	Refresh(
		ctx context.Context,
		req RefreshRequest,
	) (*LoginResponse, error)

	Logout(
		ctx context.Context,
		userID uuid.UUID,
		sessionID uuid.UUID,
		ipAddress string,
		userAgent string,
	) error

	LogoutAllSessions(
		ctx context.Context,
		userID uuid.UUID,
		ipAddress string,
		userAgent string,
	) error
}
