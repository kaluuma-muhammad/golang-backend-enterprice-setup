package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api/internal/bootstrap"
	"github.com/go-api/internal/interfaces/http/middleware"
)

func RegisterAuthorizationRoutes(router *gin.RouterGroup, container *bootstrap.Container) {

	verified := router.Group("")
	verified.Use(
		container.AuthMiddleware.RequireAuth(),
		middleware.RequireVerified(),
		container.RateLimitMiddleware.Handler("verified"),
	)

	protected := router.Group("")
	protected.Use(
		container.AuthMiddleware.RequireAuth(),
		container.RateLimitMiddleware.Handler("protected"),
	)
	{
		authorization := verified.Group("/authorization")
		{
			roles := authorization.Group("/roles")
			{
				roles.GET(
					"/",
					container.AuthorizationMiddleware.RequirePermission("roles.read"),
					container.AuthorizationHandler.GetRoles,
				)
				roles.POST(
					"/create",
					container.AuthorizationMiddleware.RequirePermission("roles.create"),
					container.AuthorizationHandler.CreateRole,
				)
				roles.PUT(
					"/update/:id",
					container.AuthorizationMiddleware.RequirePermission("roles.update"),
					container.AuthorizationHandler.UpdateRole,
				)
				roles.DELETE(
					"/delete/:id",
					container.AuthorizationMiddleware.RequirePermission("roles.delete"),
					container.AuthorizationHandler.DeleteRole,
				)
				roles.GET(
					"/get/:id",
					container.AuthorizationMiddleware.RequirePermission("roles.read"),
					container.AuthorizationHandler.GetRole,
				)
			}

			permissions := authorization.Group("/permissions")
			{
				permissions.GET(
					"/",
					container.AuthorizationMiddleware.RequirePermission("permissions.read"),
					container.AuthorizationHandler.GetPermissions,
				)

				permissions.POST(
					"/create",
					container.AuthorizationMiddleware.RequirePermission("permissions.create"),
					container.AuthorizationHandler.CreatePermission,
				)

				permissions.PUT(
					"/update/:id",
					container.AuthorizationMiddleware.RequirePermission("permissions.update"),
					container.AuthorizationHandler.UpdatePermission,
				)

				permissions.DELETE(
					"/delete/:id",
					container.AuthorizationMiddleware.RequirePermission("permissions.delete"),
					container.AuthorizationHandler.DeletePermission,
				)

				permissions.GET(
					"/get/:id",
					container.AuthorizationMiddleware.RequirePermission("permissions.read"),
					container.AuthorizationHandler.GetPermission,
				)
			}

			assignments := authorization.Group("/assignments")
			{
				assignments.POST(
					"/assign-role",
					container.AuthorizationMiddleware.RequirePermission("users.assign-role"),
					container.AuthorizationHandler.AssignRoleToUser,
				)
				assignments.POST(
					"/remove-role",
					container.AuthorizationMiddleware.RequirePermission("users.unassign-role"),
					container.AuthorizationHandler.RemoveRoleFromUser,
				)
				assignments.POST(
					"/assign-permission",
					container.AuthorizationMiddleware.RequirePermission("roles.assign-permission"),
					container.AuthorizationHandler.AssignPermissionToRole,
				)
				assignments.POST(
					"/remove-permission",
					container.AuthorizationMiddleware.RequirePermission("roles.unassign-permission"),
					container.AuthorizationHandler.RemovePermissionFromRole,
				)

				assignments.GET(
					"/list-user-roles/:userId",
					container.AuthorizationMiddleware.RequirePermission("users.read"),
					container.AuthorizationHandler.ListUserRoles,
				)

				assignments.GET(
					"/list-user-permissions/:userId",
					container.AuthorizationMiddleware.RequirePermission("users.read"),
					container.AuthorizationHandler.ListUserPermissions,
				)

				assignments.GET(
					"/list-role-permissions/:roleId",
					container.AuthorizationMiddleware.RequirePermission("roles.read"),
					container.AuthorizationHandler.ListRolePermissions,
				)
			}

			checks := authorization.Group("/checks")
			{
				checks.GET(
					"/user-has-role/:userId/:roleId",
					container.AuthorizationMiddleware.RequirePermission("users.read"),
					container.AuthorizationHandler.UserHasRole,
				)
				checks.GET(
					"/user-has-permission/:userId/:permission",
					container.AuthorizationMiddleware.RequirePermission("users.read"),
					container.AuthorizationHandler.UserHasPermission,
				)
				checks.GET(
					"/role-has-permission/:roleId/:permissionId",
					container.AuthorizationMiddleware.RequirePermission("roles.read"),
					container.AuthorizationHandler.RoleHasPermission,
				)
			}
		}
	}
}
