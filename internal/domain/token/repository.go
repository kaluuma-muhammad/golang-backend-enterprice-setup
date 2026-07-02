package token

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		token *Token,
	) error

	FindByToken(
		ctx context.Context,
		token string,
	) (*Token, error)

	MarkUsed(
		ctx context.Context,
		id uuid.UUID,
	) error

	DeleteExpired(
		ctx context.Context,
	) error
}
