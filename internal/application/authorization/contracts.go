package authorization

import (
	"context"

	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/domain/authorization/entities"
	"github.com/google/uuid"
)

type AuthorizationServiceContract interface {
	// Roles
	GetRoles(ctx context.Context, req common.PaginationRequest) ([]*entities.Role, common.PaginationResponse, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*entities.Role, error)
	GetRoleByName(ctx context.Context, name string) (*entities.Role, error)
	CreateRole(ctx context.Context, req CreateRoleRequest) (*entities.Role, error)
	UpdateRole(ctx context.Context, req UpdateRoleRequest) (*entities.Role, error)
	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	// Permissions
	GetPermissions(ctx context.Context, req common.PaginationRequest) ([]*entities.Permission, common.PaginationResponse, error)
	GetPermission(ctx context.Context, permissionID uuid.UUID) (*entities.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*entities.Permission, error)
	CreatePermission(ctx context.Context, req CreatePermissionRequest) (*entities.Permission, error)
	UpdatePermission(ctx context.Context, req UpdatePermissionRequest) (*entities.Permission, error)
	DeletePermission(ctx context.Context, permissionID uuid.UUID) error

	// Assignments
	AssignRoleToUser(ctx context.Context, req AssignRoleRequest) error
	RemoveRoleFromUser(ctx context.Context, req AssignRoleRequest) error

	AssignPermissionToRole(ctx context.Context, req AssignPermissionRequest) error
	RemovePermissionFromRole(ctx context.Context, req AssignPermissionRequest) error

	// Lists
	ListUserRoles(ctx context.Context, userID uuid.UUID) ([]*entities.Role, error)
	ListUserPermissions(ctx context.Context, userID uuid.UUID) ([]*entities.Permission, error)
	ListRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*entities.Permission, error)

	// Checks
	UserHasRole(ctx context.Context, userID, roleID uuid.UUID) (bool, error)
	UserHasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)

	RoleHasPermission(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error)

	UserHasAnyPermission(ctx context.Context, userID uuid.UUID, permissions ...string) (bool, error)
	UserHasAllPermissions(ctx context.Context, userID uuid.UUID, permissions ...string) (bool, error)

	Authorize(ctx context.Context, userID uuid.UUID, permissionName string) error
	AuthorizeAny(ctx context.Context, userID uuid.UUID, permissions ...string) error
	AuthorizeAll(ctx context.Context, userID uuid.UUID, permissions ...string) error
}
