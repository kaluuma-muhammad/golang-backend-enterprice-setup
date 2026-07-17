package authorization

import "errors"

var (
	ErrRoleNameRequired              = errors.New("role name is required")
	ErrRoleDescriptionRequired       = errors.New("role description is required")
	ErrRoleIsSystemRequired          = errors.New("role is_system is required")
	ErrPermissionNameRequired        = errors.New("permission name is required")
	ErrPermissionDescriptionRequired = errors.New("permission description is required")
	ErrPermissionResourceRequired    = errors.New("permission resource is required")
	ErrPermissionActionRequired      = errors.New("permission action is required")

	ErrRoleAlreadyExists         = errors.New("role already exists")
	ErrPermissionAlreadyExists   = errors.New("permission already exists")
	ErrRoleNotFound              = errors.New("role not found")
	ErrPermissionNotFound        = errors.New("permission not found")
	ErrUserAlreadyHasRole        = errors.New("user already has role")
	ErrRoleAlreadyHasPermission  = errors.New("role already has permission")
	ErrForbidden                 = errors.New("forbidden")
	ErrSystemRoleImmutable       = errors.New("system role cannot be modified")
	ErrSystemRoleCantBeDeleted   = errors.New("system role cannot be deleted")
	ErrUserRoleNotAssigned       = errors.New("user role not assigned")
	ErrRolePermissionNotAssigned = errors.New("role permission not assigned")

	ErrUserNotFound = errors.New("user not found")
)
