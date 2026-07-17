package authorization

import (
	"context"
	"errors"

	"github.com/go-api/internal/domain/authorization/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) AssignRoleToUser(ctx context.Context, req AssignRoleRequest) error {
	_, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	_, err = s.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRoleNotFound
		}
		return err
	}

	hasRole, err := s.assignmentRepo.UserHasRole(ctx, req.UserID, req.RoleID)
	if err != nil {
		return err
	}

	if hasRole {
		return ErrUserAlreadyHasRole
	}

	return s.assignmentRepo.AssignRoleToUser(ctx, req.UserID, req.RoleID)
}

func (s *Service) RemoveRoleFromUser(ctx context.Context, req AssignRoleRequest) error {
	_, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	role, err := s.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRoleNotFound
		}
		return err
	}

	if role.IsSystem {
		return ErrSystemRoleImmutable
	}

	hasRole, err := s.assignmentRepo.UserHasRole(ctx, req.UserID, req.RoleID)
	if err != nil {
		return err
	}

	if !hasRole {
		return ErrUserRoleNotAssigned
	}

	return s.assignmentRepo.RemoveRoleFromUser(ctx, req.UserID, req.RoleID)
}

func (s *Service) AssignPermissionToRole(ctx context.Context, req AssignPermissionRequest) error {
	_, err := s.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRoleNotFound
		}
		return err
	}

	_, err = s.permissionRepo.FindPermissionByID(ctx, req.PermissionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPermissionNotFound
		}
		return err
	}

	hasPermission, err := s.assignmentRepo.RoleHasPermission(ctx, req.RoleID, req.PermissionID)
	if err != nil {
		return err
	}

	if hasPermission {
		return ErrRoleAlreadyHasPermission
	}

	return s.assignmentRepo.AssignPermissionToRole(ctx, req.RoleID, req.PermissionID)
}

func (s *Service) RemovePermissionFromRole(ctx context.Context, req AssignPermissionRequest) error {
	role, err := s.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRoleNotFound
		}
		return err
	}

	if role.IsSystem {
		return ErrSystemRoleImmutable
	}

	_, err = s.permissionRepo.FindPermissionByID(ctx, req.PermissionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPermissionNotFound
		}
		return err
	}

	hasPermission, err := s.assignmentRepo.RoleHasPermission(ctx, req.RoleID, req.PermissionID)
	if err != nil {
		return err
	}

	if !hasPermission {
		return ErrRolePermissionNotAssigned
	}

	return s.assignmentRepo.RemovePermissionFromRole(ctx, req.RoleID, req.PermissionID)
}

func (s *Service) ListUserRoles(ctx context.Context, userID uuid.UUID) ([]*entities.Role, error) {
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return s.assignmentRepo.ListUserRoles(ctx, userID)
}

func (s *Service) ListUserPermissions(ctx context.Context, userID uuid.UUID) ([]*entities.Permission, error) {
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return s.assignmentRepo.ListUserPermissions(ctx, userID)
}

func (s *Service) ListRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*entities.Permission, error) {
	_, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	return s.assignmentRepo.ListRolePermissions(ctx, roleID)
}

func (s *Service) RoleHasPermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) (bool, error) {
	if _, err := s.requireRole(ctx, roleID); err != nil {
		return false, err
	}

	if _, err := s.requirePermission(ctx, permissionID); err != nil {
		return false, err
	}

	return s.assignmentRepo.RoleHasPermission(ctx, roleID, permissionID)
}
