package repositories

import (
	"context"
	"time"

	"github.com/go-api/internal/domain/session"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/go-api/internal/infrastructure/postgres/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SessionRepository struct {
	q *db.Queries
}

func NewSessionRepository(q *db.Queries) *SessionRepository {
	return &SessionRepository{
		q: q,
	}
}

func toSessionDomain(s db.Session) *session.Session {
	return &session.Session{
		ID:            s.ID,
		UserID:        s.UserID,
		RefreshToken:  s.RefreshToken,
		UserAgent:     types.FromPGText(s.UserAgent),
		IPAddress:     types.FromPGInet(s.IpAddress),
		DeviceName:    types.FromPGText(s.DeviceName),
		LastUsedAt:    types.FromPGTimestamp(s.LastUsedAt),
		ExpiresAt:     types.FromPGTimestamp(s.ExpiresAt),
		RevokedAt:     types.FromPGTimestampPtr(s.RevokedAt),
		RevokedReason: types.FromNullablePGText(s.RevokedReason),
		CreatedAt:     types.FromPGTimestamp(s.CreatedAt),
		UpdatedAt:     types.FromPGTimestamp(s.UpdatedAt),
	}
}

func (r *SessionRepository) Create(ctx context.Context, s *session.Session) error {
	return r.q.CreateSession(
		ctx,
		db.CreateSessionParams{

			ID:            s.ID,
			UserID:        s.UserID,
			RefreshToken:  s.RefreshToken,
			UserAgent:     types.ToPGText(s.UserAgent),
			IpAddress:     types.ToPGInet(s.IPAddress),
			DeviceName:    types.ToPGText(s.DeviceName),
			RevokedReason: types.ToNullablePGText(s.RevokedReason),
			RevokedAt:     types.ToPGTimestampPtr(s.RevokedAt),
			LastUsedAt:    types.ToPGTimestamp(s.LastUsedAt),
			ExpiresAt:     types.ToPGTimestamp(s.ExpiresAt),
			CreatedAt:     types.ToPGTimestamp(s.CreatedAt),
			UpdatedAt:     types.ToPGTimestamp(s.UpdatedAt),
		},
	)
}

func (r *SessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*session.Session, error) {

	record, err := r.q.GetSessionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toSessionDomain(record), nil
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*session.Session, error) {

	record, err := r.q.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return toSessionDomain(record), nil
}

func (r *SessionRepository) UpdateRefreshToken(ctx context.Context, id uuid.UUID, token string, expires time.Time) error {

	return r.q.UpdateSessionRefreshToken(
		ctx,
		db.UpdateSessionRefreshTokenParams{
			ID:           id,
			RefreshToken: token,
			ExpiresAt:    types.ToPGTimestamp(expires),
		},
	)
}

func (r *SessionRepository) RevokeRefreshToken(ctx context.Context, id uuid.UUID, reason string) error {
	return r.q.RevokeRefreshToken(
		ctx,
		db.RevokeRefreshTokenParams{
			ID:            id,
			RevokedReason: pgtype.Text{String: reason, Valid: reason != ""},
		},
	)
}

func (r *SessionRepository) RevokeAllSessions(ctx context.Context, userID uuid.UUID, reason string) error {
	return r.q.RevokeAllSessions(
		ctx,
		db.RevokeAllSessionsParams{
			UserID:        userID,
			RevokedReason: pgtype.Text{String: reason, Valid: reason != ""},
		},
	)
}

func (r *SessionRepository) UpdateLastUsedAt(ctx context.Context, id uuid.UUID) error {

	return r.q.UpdateSessionLastUsed(ctx, id)
}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {

	return r.q.DeleteSessionByID(ctx, id)
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {

	return r.q.DeleteSessionsByUserID(ctx, userID)
}
