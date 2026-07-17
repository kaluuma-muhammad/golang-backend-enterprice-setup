package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/domain/audit"

	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/go-api/internal/infrastructure/postgres/types"
)

type AuditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{
		pool: pool,
	}
}

func (r *AuditRepository) Create(ctx context.Context, log *audit.Log) error {
	q := common.GetQueries(ctx, r.pool)

	return q.CreateAuditLog(
		ctx,
		db.CreateAuditLogParams{
			ID:         log.ID,
			UserID:     types.ToNullableUUID(log.UserID),
			Action:     log.Action,
			EntityType: types.ToNullablePGText(log.EntityType),
			EntityID:   types.ToNullableUUID(log.EntityID),
			IpAddress:  types.ToPGInet(log.IPAddress),
			UserAgent:  types.ToPGText(log.UserAgent),
			Metadata:   types.ToJSONB(log.Metadata),
			CreatedAt:  types.ToPGTimestamp(log.CreatedAt),
			UpdatedAt:  types.ToPGTimestamp(log.UpdatedAt),
		},
	)
}

func (r *AuditRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*audit.Log, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.GetAuditLogsByUser(ctx, types.ToNullableUUID(&userID))

	if err != nil {
		return nil, err
	}

	result := make([]*audit.Log, 0, len(records))

	for _, record := range records {
		result = append(result, toAuditLogDomain(record))
	}

	return result, nil
}

func (r *AuditRepository) GetAllLogs(ctx context.Context) ([]*audit.Log, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.GetAuditLogs(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*audit.Log, 0, len(records))

	for _, record := range records {
		result = append(result, toAuditLogDomain(record))
	}

	return result, nil
}

func (r *AuditRepository) GetByUserPaginated(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*audit.Log, int64, error) {

	q := common.GetQueries(ctx, r.pool)

	records, err := q.GetAuditLogsByUserPaginated(ctx, db.GetAuditLogsByUserPaginatedParams{
		UserID: types.ToNullableUUID(&userID),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := q.CountAuditLogsByUser(ctx, types.ToNullableUUID(&userID))
	if err != nil {
		return nil, 0, err
	}

	var result []*audit.Log
	for _, r := range records {
		result = append(result, toAuditLogDomain(r))
	}

	return result, total, nil
}

func (r *AuditRepository) GetUserLogs(ctx context.Context, userID uuid.UUID) ([]*audit.Log, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.GetAuditLogsByUser(ctx, pgtype.UUID{Bytes: userID})
	if err != nil {
		return nil, err
	}

	result := make([]*audit.Log, 0, len(records))

	for _, record := range records {
		result = append(result, toAuditLogDomain(record))
	}

	return result, nil
}

func (r *AuditRepository) FindByID(ctx context.Context, id uuid.UUID) (*audit.Log, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetAuditLogByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toAuditLogDomain(record), nil
}
