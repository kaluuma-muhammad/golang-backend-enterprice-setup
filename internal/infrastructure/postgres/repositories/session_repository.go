package repositories

import (
	"context"
	"time"

	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/domain/session"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/go-api/internal/infrastructure/postgres/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		pool: pool,
	}
}

func toSessionDomain(s db.Session) *session.Session {
	return &session.Session{
		ID:            s.ID,
		UserID:        s.UserID,
		RefreshToken:  s.RefreshToken,
		UserAgent:     types.FromPGText(s.UserAgent),
		IPAddress:     types.FromPGInet(s.IpAddress),
		DeviceID:      types.FromPGText(s.DeviceID),
		Platform:      types.FromPGText(s.Platform),
		Browser:       types.FromPGText(s.Browser),
		LastSeenAt:    types.FromPGTimestamp(s.LastSeenAt),
		LastUsedAt:    types.FromPGTimestampPtr(s.LastUsedAt),
		ExpiresAt:     types.FromPGTimestamp(s.ExpiresAt),
		RevokedAt:     types.FromPGTimestampPtr(s.RevokedAt),
		RevokedReason: types.FromNullablePGText(s.RevokedReason),
		CreatedAt:     types.FromPGTimestamp(s.CreatedAt),
		UpdatedAt:     types.FromPGTimestamp(s.UpdatedAt),
	}
}

func (r *SessionRepository) Create(ctx context.Context, s *session.Session) error {
	q := common.GetQueries(ctx, r.pool)

	return q.CreateSession(
		ctx,
		db.CreateSessionParams{
			ID:            s.ID,
			UserID:        s.UserID,
			RefreshToken:  s.RefreshToken,
			UserAgent:     types.ToPGText(s.UserAgent),
			IpAddress:     types.ToPGInet(s.IPAddress),
			DeviceID:      types.ToPGText(s.DeviceID),
			Platform:      types.ToPGText(s.Platform),
			Browser:       types.ToPGText(s.Browser),
			LastSeenAt:    types.ToPGTimestamp(s.LastSeenAt),
			LastUsedAt:    types.ToPGTimestampPtr(s.LastUsedAt),
			ExpiresAt:     types.ToPGTimestamp(s.ExpiresAt),
			RevokedAt:     types.ToPGTimestampPtr(s.RevokedAt),
			RevokedReason: types.ToNullablePGText(s.RevokedReason),
			CreatedAt:     types.ToPGTimestamp(s.CreatedAt),
			UpdatedAt:     types.ToPGTimestamp(s.UpdatedAt),
		},
	)
}

func (r *SessionRepository) GetSessionsByUserPaginated(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*session.Session, int64, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.GetSessionsByUserPaginated(ctx, db.GetSessionsByUserPaginatedParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := q.CountSessionsByUser(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	var sessions []*session.Session
	for _, record := range records {
		sessions = append(sessions, toSessionDomain(record))
	}

	return sessions, total, nil
}

func (r *SessionRepository) GetCurrentSessions(ctx context.Context, userID uuid.UUID) ([]*session.Session, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.GetCurrentSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var sessions []*session.Session
	for _, record := range records {
		sessions = append(sessions, toSessionDomain(record))
	}

	return sessions, nil
}

func (r *SessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*session.Session, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetSessionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toSessionDomain(record), nil
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*session.Session, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return toSessionDomain(record), nil
}

func (r *SessionRepository) GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*session.Session, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.GetSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sessions := make([]*session.Session, len(records))
	for i, record := range records {
		sessions[i] = toSessionDomain(record)
	}

	return sessions, nil
}

func (r *SessionRepository) UpdateRefreshToken(ctx context.Context, id uuid.UUID, token string, expires time.Time) error {
	q := common.GetQueries(ctx, r.pool)

	return q.UpdateSessionRefreshToken(
		ctx,
		db.UpdateSessionRefreshTokenParams{
			ID:           id,
			RefreshToken: token,
			ExpiresAt:    types.ToPGTimestamp(expires),
		},
	)
}

func (r *SessionRepository) RevokeSession(ctx context.Context, id uuid.UUID, reason string) error {
	q := common.GetQueries(ctx, r.pool)
	return q.RevokeSession(
		ctx,
		db.RevokeSessionParams{
			ID:            id,
			RevokedReason: pgtype.Text{String: reason, Valid: reason != ""},
		},
	)
}

func (r *SessionRepository) RevokeAllSessions(ctx context.Context, userID uuid.UUID, reason string) error {
	q := common.GetQueries(ctx, r.pool)
	return q.RevokeAllSessions(
		ctx,
		db.RevokeAllSessionsParams{
			UserID:        userID,
			RevokedReason: pgtype.Text{String: reason, Valid: reason != ""},
		},
	)
}

func (r *SessionRepository) UpdateSessionActivity(ctx context.Context, id uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)
	return q.UpdateSessionActivity(ctx, id)
}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)
	return q.DeleteSessionByID(ctx, id)
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)
	return q.DeleteSessionsByUserID(ctx, userID)
}
