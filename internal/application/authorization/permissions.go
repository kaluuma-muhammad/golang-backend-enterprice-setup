package authorization

import (
	"context"
	"errors"
	"strings"

	common "github.com/go-api/internal/application/common"
	"github.com/go-api/internal/domain/authorization/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) GetPermissions(ctx context.Context, req common.PaginationRequest) ([]*entities.Permission, common.PaginationResponse, error) {
	req = req.Normalize()

	data, total, err := s.permissionRepo.GetPermissions(ctx, req.PageSize, req.Offset())
	if err != nil {
		return nil, common.PaginationResponse{}, err
	}

	return data, common.PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
	}, nil
}

func (s *Service) CreatePermission(ctx context.Context, req CreatePermissionRequest) (*entities.Permission, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Resource = strings.TrimSpace(strings.ToLower(req.Resource))
	req.Action = strings.TrimSpace(strings.ToLower(req.Action))
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return nil, ErrPermissionNameRequired
	}

	if req.Resource == "" {
		return nil, ErrPermissionResourceRequired
	}

	if req.Action == "" {
		return nil, ErrPermissionActionRequired
	}

	existing, err := s.permissionRepo.FindPermissionByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, ErrPermissionAlreadyExists
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	permission := &entities.Permission{
		ID:          uuid.New(),
		Name:        req.Name,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: &req.Description,
	}

	return s.permissionRepo.CreatePermission(ctx, permission)
}

func (s *Service) UpdatePermission(ctx context.Context, req UpdatePermissionRequest) (*entities.Permission, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Resource = strings.TrimSpace(strings.ToLower(req.Resource))
	req.Action = strings.TrimSpace(strings.ToLower(req.Action))
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return nil, ErrPermissionNameRequired
	}

	if req.Resource == "" {
		return nil, ErrPermissionResourceRequired
	}

	if req.Action == "" {
		return nil, ErrPermissionActionRequired
	}

	permission, err := s.permissionRepo.FindPermissionByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}

	existing, err := s.permissionRepo.FindPermissionByName(ctx, req.Name)
	if err == nil && existing != nil && existing.ID != permission.ID {
		return nil, ErrPermissionAlreadyExists
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	permission.Name = req.Name
	permission.Resource = req.Resource
	permission.Action = req.Action
	permission.Description = &req.Description

	return s.permissionRepo.UpdatePermission(ctx, permission)
}

func (s *Service) GetPermission(ctx context.Context, permissionID uuid.UUID) (*entities.Permission, error) {
	permission, err := s.permissionRepo.FindPermissionByID(ctx, permissionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}

	return permission, nil
}

func (s *Service) GetPermissionByName(ctx context.Context, name string) (*entities.Permission, error) {
	permission, err := s.permissionRepo.FindPermissionByName(ctx, strings.TrimSpace(name))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}

	return permission, nil
}

func (s *Service) DeletePermission(ctx context.Context, permissionID uuid.UUID) error {
	_, err := s.permissionRepo.FindPermissionByID(ctx, permissionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPermissionNotFound
		}
		return err
	}

	return s.permissionRepo.DeletePermission(ctx, permissionID)
}
