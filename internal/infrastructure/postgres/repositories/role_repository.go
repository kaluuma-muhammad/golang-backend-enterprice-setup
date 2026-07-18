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

type RoleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{
		pool: pool,
	}
}

func (r *RoleRepository) CreateRole(ctx context.Context, role *entities.Role) (*entities.Role, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.CreateRole(
		ctx,
		db.CreateRoleParams{
			ID:          role.ID,
			Name:        role.Name,
			Description: types.ToNullablePGText(role.Description),
			IsSystem:    role.IsSystem,
		},
	)

	if err != nil {
		return nil, err
	}

	return toRoleDomain(record), nil
}

func (r *RoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toRoleDomain(record), nil
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*entities.Role, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetRoleByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return toRoleDomain(record), nil
}

func (r *RoleRepository) GetRoles(ctx context.Context, limit int, offset int) ([]*entities.Role, int64, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.GetPaginatedRoles(
		ctx,
		db.GetPaginatedRolesParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		},
	)

	if err != nil {
		return nil, 0, err
	}

	total, err := q.CountRoles(ctx)
	if err != nil {
		return nil, 0, err
	}

	var results = make([]*entities.Role, len(records))
	for i, record := range records {
		results[i] = toRoleDomain(record)
	}

	return results, total, nil
}

func (r *RoleRepository) UpdateRole(ctx context.Context, role *entities.Role) (*entities.Role, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.UpdateRole(
		ctx,
		db.UpdateRoleParams{
			ID:          role.ID,
			Name:        role.Name,
			Description: types.ToNullablePGText(role.Description),
		},
	)

	if err != nil {
		return nil, err
	}

	return toRoleDomain(record), nil
}

func (r *RoleRepository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)

	return q.DeleteRole(ctx, id)
}
