package auth

import (
	"context"

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
