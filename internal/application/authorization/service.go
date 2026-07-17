package authorization

import (
	authorization "github.com/go-api/internal/domain/authorization/repositories"
	userRepo "github.com/go-api/internal/domain/user"
)

type Service struct {
	roleRepo       authorization.RoleRepository
	permissionRepo authorization.PermissionRepository
	assignmentRepo authorization.AssignmentRepository
	userRepo       userRepo.Repository
}

func NewService(
	roleRepo authorization.RoleRepository,
	permissionRepo authorization.PermissionRepository,
	assignmentRepo authorization.AssignmentRepository,
	userRepo userRepo.Repository,
) *Service {
	return &Service{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		assignmentRepo: assignmentRepo,
		userRepo:       userRepo,
	}
}
