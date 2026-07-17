package authorization

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) UserHasRole(ctx context.Context, userID, roleID uuid.UUID) (bool, error) {
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		return false, err
	}

	if _, err := s.requireRole(ctx, roleID); err != nil {
		return false, err
	}

	return s.assignmentRepo.UserHasRole(
		ctx,
		userID,
		roleID,
	)
}

func (s *Service) UserHasPermission(ctx context.Context, userID uuid.UUID, permissionName string) (bool, error) {
	if _, err := s.requireUser(ctx, userID); err != nil {
		return false, err
	}

	return s.assignmentRepo.UserHasPermission(
		ctx,
		userID,
		permissionName,
	)
}

func (s *Service) Authorize(ctx context.Context, userID uuid.UUID, permissionName string) error {
	allowed, err := s.UserHasPermission(ctx, userID, permissionName)
	if err != nil {
		return err
	}

	if !allowed {
		return ErrForbidden
	}

	return nil
}

func (s *Service) UserHasAnyPermission(ctx context.Context, userID uuid.UUID, permissions ...string) (bool, error) {
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		return false, err
	}

	for _, permission := range permissions {
		ok, err := s.assignmentRepo.UserHasPermission(ctx, userID, permission)
		if err != nil {
			return false, err
		}

		if ok {
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) UserHasAllPermissions(ctx context.Context, userID uuid.UUID, permissions ...string) (bool, error) {
	if _, err := s.requireUser(ctx, userID); err != nil {
		return false, err
	}

	if len(permissions) == 0 {
		return false, nil
	}

	for _, permission := range permissions {
		ok, err := s.assignmentRepo.UserHasPermission(ctx, userID, permission)
		if err != nil {
			return false, err
		}

		if !ok {
			return false, nil
		}
	}

	return true, nil
}

func (s *Service) AuthorizeAny(ctx context.Context, userID uuid.UUID, permissions ...string) error {
	allowed, err := s.UserHasAnyPermission(ctx, userID, permissions...)
	if err != nil {
		return err
	}

	if !allowed {
		return ErrForbidden
	}

	return nil
}

func (s *Service) AuthorizeAll(ctx context.Context, userID uuid.UUID, permissions ...string) error {
	allowed, err := s.UserHasAllPermissions(ctx, userID, permissions...)
	if err != nil {
		return err
	}

	if !allowed {
		return ErrForbidden
	}

	return nil
}
