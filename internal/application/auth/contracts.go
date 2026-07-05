package auth

import (
	"context"

	backendUser "github.com/go-api/internal/domain/user"
	"github.com/google/uuid"
)

type ServiceContract interface {
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

	VerifyAccount(
		ctx context.Context,
		userID uuid.UUID,
		req VerifyAccountRequest,
	) error

	ForgotPassword(
		ctx context.Context,
		req ForgotPasswordRequest,
	) error

	ResendVerification(
		ctx context.Context,
		req ResendVerificationRequest,
	) error

	VerifyEmail(
		ctx context.Context,
		req VerifyEmailRequest,
		userAgent string,
		ipAddress string,
	) (*LoginResponse, error)

	ResetPassword(
		ctx context.Context,
		user *backendUser.User,
		req ResetPasswordRequest,
		userAgent string,
		ipAddress string,
	) (*LoginResponse, error)

	Refresh(
		ctx context.Context,
		req RefreshRequest,
	) (*LoginResponse, error)

	Logout(
		ctx context.Context,
		sessionID uuid.UUID,
	) error

	LogoutAllSessions(
		ctx context.Context,
		userID uuid.UUID,
	) error
}
