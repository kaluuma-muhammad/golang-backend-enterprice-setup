package authorization

import (
	"context"
	"errors"

	"github.com/go-api/internal/domain/authorization/entities"
	"github.com/go-api/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) requireUser(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func (s *Service) requireRole(ctx context.Context, roleID uuid.UUID) (*entities.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	return role, nil
}

func (s *Service) requirePermission(ctx context.Context, permissionID uuid.UUID) (*entities.Permission, error) {
	permission, err := s.permissionRepo.FindPermissionByID(ctx, permissionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}

	return permission, nil
}

func (s *Service) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]*entities.Permission, error) {
	return s.ListUserPermissions(ctx, userID)
}
