package repositories

import (
	"context"

	"github.com/go-api/internal/domain/authorization/entities"
	"github.com/google/uuid"
)

type RoleRepository interface {
	CreateRole(
		ctx context.Context,
		role *entities.Role,
	) (*entities.Role, error)

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entities.Role, error)

	FindByName(
		ctx context.Context,
		name string,
	) (*entities.Role, error)

	GetRoles(
		ctx context.Context,
		limit int,
		offset int,
	) ([]*entities.Role, int64, error)

	UpdateRole(
		ctx context.Context,
		role *entities.Role,
	) (*entities.Role, error)

	DeleteRole(
		ctx context.Context,
		id uuid.UUID,
	) error
}
