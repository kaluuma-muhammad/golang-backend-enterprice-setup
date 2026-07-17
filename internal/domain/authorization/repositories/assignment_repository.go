package repositories

import (
	"context"

	"github.com/go-api/internal/domain/authorization/entities"
	"github.com/google/uuid"
)

type AssignmentRepository interface {
	AssignRoleToUser(
		ctx context.Context,
		userID uuid.UUID,
		roleID uuid.UUID,
	) error

	RemoveRoleFromUser(
		ctx context.Context,
		userID uuid.UUID,
		roleID uuid.UUID,
	) error

	AssignPermissionToRole(
		ctx context.Context,
		roleID uuid.UUID,
		permissionID uuid.UUID,
	) error

	RemovePermissionFromRole(
		ctx context.Context,
		roleID uuid.UUID,
		permissionID uuid.UUID,
	) error

	ListUserRoles(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*entities.Role, error)

	ListRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
	) ([]*entities.Permission, error)

	ListUserPermissions(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*entities.Permission, error)

	UserHasRole(
		ctx context.Context,
		userID uuid.UUID,
		roleID uuid.UUID,
	) (bool, error)

	RoleHasPermission(
		ctx context.Context,
		roleID uuid.UUID,
		permissionID uuid.UUID,
	) (bool, error)

	UserHasPermission(
		ctx context.Context,
		userID uuid.UUID,
		permission string,
	) (bool, error)
}
