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

	FindByTokenAndType(
		ctx context.Context,
		token string,
		tokenType Type,
	) (*Token, error)

	FindByTokenAndTypeAndUser(
		ctx context.Context,
		token string,
		tokenType Type,
		userID uuid.UUID,
	) (*Token, error)

	FindByUserAndType(
		ctx context.Context,
		userID uuid.UUID,
		tokenType Type,
	) (*Token, error)

	MarkUsed(
		ctx context.Context,
		id uuid.UUID,
	) error

	DeleteByUserAndType(
		ctx context.Context,
		userID uuid.UUID,
		tokenType Type,
	) error

	DeleteExpired(
		ctx context.Context,
	) error
}
