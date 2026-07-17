package repositories

import (
	"context"

	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/domain/authorization/entities"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/go-api/internal/infrastructure/postgres/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionRepository struct {
	pool *pgxpool.Pool
}

func NewPermissionRepository(pool *pgxpool.Pool) *PermissionRepository {
	return &PermissionRepository{
		pool: pool,
	}
}

func (r *PermissionRepository) CreatePermission(ctx context.Context, permission *entities.Permission) (*entities.Permission, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.CreatePermission(
		ctx,
		db.CreatePermissionParams{
			Resource:    permission.Resource,
			Action:      permission.Action,
			Name:        permission.Name,
			Description: types.ToNullablePGText(permission.Description),
		},
	)

	if err != nil {
		return nil, err
	}

	return toPermissionDomain(record), nil
}

func (r *PermissionRepository) FindPermissionByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetPermissionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toPermissionDomain(record), nil
}

func (r *PermissionRepository) FindPermissionByName(ctx context.Context, name string) (*entities.Permission, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetPermissionByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return toPermissionDomain(record), nil
}

func (r *PermissionRepository) GetPermissions(ctx context.Context, limit int, offset int) ([]*entities.Permission, int64, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.GetPaginatedPermissions(
		ctx,
		db.GetPaginatedPermissionsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		},
	)

	if err != nil {
		return nil, 0, err
	}

	total, err := q.CountPermissions(ctx)
	if err != nil {
		return nil, 0, err
	}

	permissions := make([]*entities.Permission, len(records))
	for i, record := range records {
		permissions[i] = toPermissionDomain(record)
	}

	return permissions, total, nil
}

func (r *PermissionRepository) UpdatePermission(ctx context.Context, permission *entities.Permission) (*entities.Permission, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.UpdatePermission(
		ctx,
		db.UpdatePermissionParams{
			ID:          permission.ID,
			Resource:    permission.Resource,
			Action:      permission.Action,
			Name:        permission.Name,
			Description: types.ToNullablePGText(permission.Description),
		},
	)

	if err != nil {
		return nil, err
	}

	return toPermissionDomain(record), nil
}

func (r *PermissionRepository) DeletePermission(ctx context.Context, id uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)

	return q.DeletePermission(ctx, id)
}
