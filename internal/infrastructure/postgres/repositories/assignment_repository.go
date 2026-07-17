package repositories

import (
	"context"

	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/domain/authorization/entities"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssignmentRepository struct {
	pool *pgxpool.Pool
}

func NewAssignmentRepository(pool *pgxpool.Pool) *AssignmentRepository {
	return &AssignmentRepository{
		pool: pool,
	}
}

func (r *AssignmentRepository) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)

	return q.AssignRoleToUser(
		ctx,
		db.AssignRoleToUserParams{
			UserID: userID,
			RoleID: roleID,
		},
	)
}

func (r *AssignmentRepository) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)

	return q.RemoveRoleFromUser(
		ctx,
		db.RemoveRoleFromUserParams{
			UserID: userID,
			RoleID: roleID,
		},
	)
}

func (r *AssignmentRepository) AssignPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)

	return q.AssignPermissionToRole(
		ctx,
		db.AssignPermissionToRoleParams{
			RoleID:       roleID,
			PermissionID: permissionID,
		},
	)
}

func (r *AssignmentRepository) RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)

	return q.RemovePermissionFromRole(
		ctx,
		db.RemovePermissionFromRoleParams{
			RoleID:       roleID,
			PermissionID: permissionID,
		},
	)
}

func (r *AssignmentRepository) ListUserRoles(ctx context.Context, userID uuid.UUID) ([]*entities.Role, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.ListRolesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles := toRoleDomains(records)

	return roles, nil
}

func (r *AssignmentRepository) ListRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*entities.Permission, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.ListPermissionsByRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	permissions := toPermissionDomains(records)

	return permissions, nil
}

func (r *AssignmentRepository) ListUserPermissions(ctx context.Context, userID uuid.UUID) ([]*entities.Permission, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.ListUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissions := toPermissionDomains(records)

	return permissions, nil
}

func (r *AssignmentRepository) ListPermissionRoles(ctx context.Context, permissionID uuid.UUID) ([]*entities.Role, error) {
	q := common.GetQueries(ctx, r.pool)

	records, err := q.ListRolesByPermission(ctx, permissionID)
	if err != nil {
		return nil, err
	}

	roles := toRoleDomains(records)

	return roles, nil
}

func (r *AssignmentRepository) UserHasRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (bool, error) {
	q := common.GetQueries(ctx, r.pool)

	return q.UserHasRole(
		ctx,
		db.UserHasRoleParams{
			UserID: userID,
			RoleID: roleID,
		},
	)
}

func (r *AssignmentRepository) RoleHasPermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) (bool, error) {
	q := common.GetQueries(ctx, r.pool)

	return q.RoleHasPermission(
		ctx,
		db.RoleHasPermissionParams{
			RoleID:       roleID,
			PermissionID: permissionID,
		},
	)
}

func (r *AssignmentRepository) UserHasPermission(ctx context.Context, userID uuid.UUID, permissionName string) (bool, error) {
	q := common.GetQueries(ctx, r.pool)

	return q.UserHasPermission(
		ctx,
		db.UserHasPermissionParams{
			UserID: userID,
			Name:   permissionName,
		},
	)
}

func (r *AssignmentRepository) HasPermission(ctx context.Context, roleName string, resource string, action string) (bool, error) {
	q := common.GetQueries(ctx, r.pool)

	return q.HasPermission(
		ctx,
		db.HasPermissionParams{
			Name:     roleName,
			Resource: resource,
			Action:   action,
		},
	)
}
