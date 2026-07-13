package audit

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		log *Log,
	) error

	GetAllLogs(
		ctx context.Context,
	) ([]*Log, error)

	GetUserLogs(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*Log, error)

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*Log, error)

	GetByUserPaginated(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
		offset int,
	) ([]*Log, int64, error)
}
