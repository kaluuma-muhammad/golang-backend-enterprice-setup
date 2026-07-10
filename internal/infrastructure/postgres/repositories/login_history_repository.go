package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/domain/loginhistory"

	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/go-api/internal/infrastructure/postgres/types"
)

type LoginHistoryRepository struct {
	pool *pgxpool.Pool
}

func NewLoginHistoryRepository(pool *pgxpool.Pool) *LoginHistoryRepository {
	return &LoginHistoryRepository{
		pool: pool,
	}
}

func toLoginHistoryDomain(h db.LoginHistory) *loginhistory.LoginHistory {
	return &loginhistory.LoginHistory{
		ID:            h.ID,
		UserID:        types.FromNullableUUID(h.UserID),
		SessionID:     types.FromNullableUUID(h.SessionID),
		Status:        loginhistory.Status(h.Status),
		IPAddress:     types.FromPGInet(h.IpAddress),
		UserAgent:     types.FromPGText(h.UserAgent),
		DeviceName:    types.FromPGText(h.DeviceName),
		FailureReason: types.FromNullablePGText(h.FailureReason),
		LoginAt:       types.FromPGTimestamp(h.LoginAt),
		LogoutAt:      types.FromPGTimestampPtr(h.LogoutAt),
		CreatedAt:     types.FromPGTimestamp(h.CreatedAt),
		UpdatedAt:     types.FromPGTimestamp(h.UpdatedAt),
	}
}

func (r *LoginHistoryRepository) Create(ctx context.Context, h *loginhistory.LoginHistory) error {
	q := common.GetQueries(ctx, r.pool)

	return q.CreateLoginHistory(
		ctx,
		db.CreateLoginHistoryParams{
			ID:            h.ID,
			UserID:        types.ToNullableUUID(h.UserID),
			SessionID:     types.ToNullableUUID(h.SessionID),
			Status:        db.LoginStatus(h.Status),
			IpAddress:     types.ToPGInet(h.IPAddress),
			UserAgent:     types.ToPGText(h.UserAgent),
			DeviceName:    types.ToPGText(h.DeviceName),
			FailureReason: types.ToNullablePGText(h.FailureReason),
			LoginAt:       types.ToPGTimestamp(h.LoginAt),
			LogoutAt:      types.ToPGTimestampPtr(h.LogoutAt),
			CreatedAt:     types.ToPGTimestamp(h.CreatedAt),
			UpdatedAt:     types.ToPGTimestamp(h.UpdatedAt),
		},
	)
}

func (r *LoginHistoryRepository) GetLoginHistoryByUser(ctx context.Context, userID uuid.UUID) ([]*loginhistory.LoginHistory, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.GetLoginHistoryByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	var histories []*loginhistory.LoginHistory
	for _, record := range records {
		histories = append(histories, toLoginHistoryDomain(record))
	}
	return histories, nil
}

func (r *LoginHistoryRepository) GetLatestSuccessfulLogin(ctx context.Context, userID uuid.UUID) (*loginhistory.LoginHistory, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetLatestSuccessfulLogin(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	return toLoginHistoryDomain(record), nil
}

func (r *LoginHistoryRepository) MarkLoginHistoryLogout(ctx context.Context, sessionID uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)

	return q.MarkLoginHistoryLogout(ctx, types.ToNullableUUID(&sessionID))
}
