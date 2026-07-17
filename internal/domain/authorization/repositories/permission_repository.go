package repositories

import (
	"context"

	"github.com/go-api/internal/domain/authorization/entities"
	"github.com/google/uuid"
)

type PermissionRepository interface {
	CreatePermission(
		ctx context.Context,
		permission *entities.Permission,
	) (*entities.Permission, error)

	FindPermissionByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entities.Permission, error)

	FindPermissionByName(
		ctx context.Context,
		name string,
	) (*entities.Permission, error)

	GetPermissions(
		ctx context.Context,
		limit int,
		offset int,
	) ([]*entities.Permission, int64, error)

	UpdatePermission(
		ctx context.Context,
		permission *entities.Permission,
	) (*entities.Permission, error)

	DeletePermission(
		ctx context.Context,
		id uuid.UUID,
	) error
}
