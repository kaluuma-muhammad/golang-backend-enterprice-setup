package seed

import (
	"context"
	"errors"

	appAuthorization "github.com/go-api/internal/application/authorization"
)

var permissions = []appAuthorization.CreatePermissionRequest{
	// Roles
	{
		Name:        "roles.read",
		Resource:    "roles",
		Action:      "read",
		Description: "View roles",
	},
	{
		Name:        "roles.create",
		Resource:    "roles",
		Action:      "create",
		Description: "Create roles",
	},
	{
		Name:        "roles.update",
		Resource:    "roles",
		Action:      "update",
		Description: "Update roles",
	},
	{
		Name:        "roles.delete",
		Resource:    "roles",
		Action:      "delete",
		Description: "Delete roles",
	},

	// Permissions
	{
		Name:        "permissions.read",
		Resource:    "permissions",
		Action:      "read",
		Description: "View permissions",
	},
	{
		Name:        "permissions.create",
		Resource:    "permissions",
		Action:      "create",
		Description: "Create permissions",
	},
	{
		Name:        "permissions.update",
		Resource:    "permissions",
		Action:      "update",
		Description: "Update permissions",
	},
	{
		Name:        "permissions.delete",
		Resource:    "permissions",
		Action:      "delete",
		Description: "Delete permissions",
	},

	// Assignments
	{
		Name:        "assignments.read",
		Resource:    "assignments",
		Action:      "read",
		Description: "View role assignments",
	},
	{
		Name:        "assignments.manage",
		Resource:    "assignments",
		Action:      "manage",
		Description: "Manage role assignments",
	},

	{
		Name:        "assignments.view-user-roles",
		Resource:    "assignments",
		Action:      "view-user-roles",
		Description: "View user roles",
	},
	{
		Name:        "assignments.view-user-permissions",
		Resource:    "assignments",
		Action:      "view-user-permissions",
		Description: "View user permissions",
	},
	{
		Name:        "assignments.view-role-permissions",
		Resource:    "assignments",
		Action:      "view-role-permissions",
		Description: "View role permissions",
	},

	// Authorization Checks
	{
		Name:        "checks.read",
		Resource:    "checks",
		Action:      "read",
		Description: "Run authorization checks",
	},

	// Users
	{
		Name:        "users.read",
		Resource:    "users",
		Action:      "read",
		Description: "View users",
	},
	{
		Name:        "users.update",
		Resource:    "users",
		Action:      "update",
		Description: "Update users",
	},
	{
		Name:        "users.delete",
		Resource:    "users",
		Action:      "delete",
		Description: "Delete users",
	},

	// audit logs
	{
		Name:        "audit-logs.view",
		Resource:    "audit-logs",
		Action:      "view",
		Description: "View audit logs",
	},
	{
		Name:        "audit-logs.export",
		Resource:    "audit-logs",
		Action:      "export",
		Description: "Export audit logs",
	},

	// login history
	{
		Name:        "login-history.view",
		Resource:    "login-history",
		Action:      "view",
		Description: "View login history",
	},

	{
		Name:        "login-history.export",
		Resource:    "login-history",
		Action:      "export",
		Description: "Export login history",
	},

	{
		Name:        "login-history.delete",
		Resource:    "login-history",
		Action:      "delete",
		Description: "Delete login history",
	},

	// session
	{
		Name:        "sessions.view",
		Resource:    "sessions",
		Action:      "view",
		Description: "View sessions",
	},
	{
		Name:        "sessions.view-active-sessions",
		Resource:    "sessions",
		Action:      "view-active-sessions",
		Description: "View active sessions",
	},
	{
		Name:        "sessions.view-expired-sessions",
		Resource:    "sessions",
		Action:      "view-expired-sessions",
		Description: "View expired sessions",
	},
}

func SeedPermissions(ctx context.Context, service *appAuthorization.Service) error {
	for _, permission := range permissions {
		_, err := service.GetPermissionByName(ctx, permission.Name)
		if err == nil {
			continue
		}

		_, err = service.CreatePermission(ctx, permission)

		if err != nil &&
			!errors.Is(err, appAuthorization.ErrPermissionAlreadyExists) {
			return err
		}
	}

	return nil
}
