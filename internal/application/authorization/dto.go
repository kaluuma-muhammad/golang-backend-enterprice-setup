package authorization

import "github.com/google/uuid"

type CreateRoleRequest struct {
	Name        string
	Description string
	IsSystem    bool
}

type UpdateRoleRequest struct {
	ID          uuid.UUID
	Name        string
	Description string
}

type CreatePermissionRequest struct {
	Resource    string
	Action      string
	Name        string
	Description string
}

type UpdatePermissionRequest struct {
	ID          uuid.UUID
	Resource    string
	Action      string
	Name        string
	Description string
}

type AssignRoleRequest struct {
	UserID uuid.UUID
	RoleID uuid.UUID
}

type AssignPermissionRequest struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}
