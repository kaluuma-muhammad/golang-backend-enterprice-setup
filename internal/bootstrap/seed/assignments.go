package seed

import (
	"context"
	"errors"

	appAuthorization "github.com/go-api/internal/application/authorization"
)

var rolePermissions = map[string][]string{
	"Super Admin": {
		// Roles
		"roles.read",
		"roles.create",
		"roles.update",
		"roles.delete",

		// Permissions
		"permissions.read",
		"permissions.create",
		"permissions.update",
		"permissions.delete",

		// Assignments
		"assignments.read",
		"assignments.manage",
		"assignments.view-user-roles",
		"assignments.view-user-permissions",
		"assignments.view-role-permissions",

		// Authorization checks
		"checks.read",

		// Users
		"users.read",
		"users.update",
		"users.delete",

		// Audit logs
		"audit-logs.view",
		"audit-logs.export",

		// Login history
		"login-history.view",
		"login-history.export",
		"login-history.delete",

		// Sessions
		"sessions.view",
		"sessions.view-active-sessions",
		"sessions.view-expired-sessions",
	},

	"Admin": {
		// Roles
		"roles.read",

		// Permissions
		"permissions.read",

		// Assignments
		"assignments.manage",
		"assignments.read",
		"assignments.view-user-roles",
		"assignments.view-user-permissions",
		"assignments.view-role-permissions",

		// Authorization checks
		"checks.read",

		// Users
		"users.read",
		"users.update",

		// Audit logs
		"audit-logs.view",

		// Login history
		"login-history.view",

		// Sessions
		"sessions.view",
		"sessions.view-active-sessions",
	},

	"User": {
		// Intentionally empty for now.
		// Self-service permissions can be added later, e.g.:
		// profile.update
		// password.change
	},
}

func SeedAssignments(ctx context.Context, service *appAuthorization.Service) error {
	for roleName, permissions := range rolePermissions {
		role, err := service.GetRoleByName(ctx, roleName)
		if err != nil {
			return err
		}

		for _, permissionName := range permissions {
			permission, err := service.GetPermissionByName(ctx, permissionName)
			if err != nil {
				return err
			}

			hasPermission, err := service.RoleHasPermission(ctx, role.ID, permission.ID)
			if err != nil {
				return err
			}

			if hasPermission {
				continue
			}

			err = service.AssignPermissionToRole(
				ctx,
				appAuthorization.AssignPermissionRequest{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				},
			)

			if err != nil &&
				!errors.Is(err, appAuthorization.ErrRoleAlreadyHasPermission) {
				return err
			}
		}
	}

	return nil
}
