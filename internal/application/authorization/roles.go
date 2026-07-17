package authorization

import (
	"context"
	"log"
	"strings"

	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/domain/authorization/entities"
	"github.com/google/uuid"
)

func (s *Service) GetRoles(ctx context.Context, req common.PaginationRequest) ([]*entities.Role, common.PaginationResponse, error) {
	req = req.Normalize()

	log.Printf("GetRoles: %v", req)

	data, total, err := s.roleRepo.GetRoles(ctx, req.PageSize, req.Offset())
	if err != nil {
		return nil, common.PaginationResponse{}, err
	}

	pagination := common.PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
	}

	return data, pagination, nil
}

func (s *Service) CreateRole(ctx context.Context, req CreateRoleRequest) (*entities.Role, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return nil, ErrRoleNameRequired
	}

	existing, err := s.roleRepo.FindByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, ErrRoleAlreadyExists
	}

	role := &entities.Role{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: &req.Description,
		IsSystem:    req.IsSystem,
	}

	return s.roleRepo.CreateRole(ctx, role)
}

func (s *Service) UpdateRole(ctx context.Context, req UpdateRoleRequest) (*entities.Role, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return nil, ErrRoleNameRequired
	}

	role, err := s.roleRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	if role.IsSystem {
		return nil, ErrSystemRoleImmutable
	}

	existing, err := s.roleRepo.FindByName(ctx, req.Name)
	if err == nil && existing != nil && existing.ID != req.ID {
		return nil, ErrRoleAlreadyExists
	}

	role.Name = req.Name
	role.Description = &req.Description

	return s.roleRepo.UpdateRole(ctx, role)
}

func (s *Service) GetRole(ctx context.Context, roleID uuid.UUID) (*entities.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	return role, nil
}

func (s *Service) GetRoleByName(ctx context.Context, name string) (*entities.Role, error) {
	role, err := s.roleRepo.FindByName(ctx, strings.TrimSpace(name))
	if err != nil {
		return nil, ErrRoleNotFound
	}

	return role, nil
}

func (s *Service) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return ErrRoleNotFound
	}

	if role.IsSystem {
		return ErrSystemRoleImmutable
	}

	return s.roleRepo.DeleteRole(ctx, roleID)
}
