package loginhistory

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		history *LoginHistory,
	) error

	GetLoginHistoryByUser(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*LoginHistory, error)

	MarkLoginHistoryLogout(
		ctx context.Context,
		sessionID uuid.UUID,
	) error

	GetLatestSuccessfulLogin(
		ctx context.Context,
		userID uuid.UUID,
	) (*LoginHistory, error)
}
