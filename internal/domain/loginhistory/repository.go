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

	GetByUserPaginated(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
		offset int,
	) ([]*LoginHistory, int64, error)

	MarkLoginHistoryLogout(
		ctx context.Context,
		sessionID uuid.UUID,
	) error

	GetLatestSuccessfulLogin(
		ctx context.Context,
		userID uuid.UUID,
	) (*LoginHistory, error)
}
