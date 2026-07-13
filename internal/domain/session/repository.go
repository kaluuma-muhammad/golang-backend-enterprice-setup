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

	GetSessionsByUserPaginated(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
		offset int,
	) ([]*Session, int64, error)

	GetCurrentSessions(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*Session, error)

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*Session, error)

	GetSessionsByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*Session, error)

	UpdateRefreshToken(
		ctx context.Context,
		sessionID uuid.UUID,
		token string,
		expiresAt time.Time,
	) error

	RevokeSession(
		ctx context.Context,
		sessionID uuid.UUID,
		reason string,
	) error

	RevokeAllSessions(
		ctx context.Context,
		userID uuid.UUID,
		reason string,
	) error

	UpdateSessionActivity(
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
