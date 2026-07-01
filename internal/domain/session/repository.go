package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		session *Session,
	) error

	FindByToken(
		ctx context.Context,
		token string,
	) (*Session, error)

	UpdateRefreshToken(
		ctx context.Context,
		sessionID uuid.UUID,
		token string,
		expiresAt time.Time,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*Session, error)

	RevokeRefreshToken(
		ctx context.Context,
		sessionID uuid.UUID,
		reason string,
	) error

	RevokeAllSessions(
		ctx context.Context,
		userID uuid.UUID,
		reason string,
	) error

	UpdateLastUsedAt(
		ctx context.Context,
		sessionID uuid.UUID,
	) error

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error

	DeleteByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) error
}
