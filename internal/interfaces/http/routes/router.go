package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/go-api/internal/bootstrap"
	"github.com/go-api/internal/interfaces/http/handlers"
	"github.com/go-api/internal/interfaces/http/middleware"
)

func SetupRouter(logger *zap.Logger, container *bootstrap.Container) *gin.Engine {

	router := gin.New()

	router.Use(
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Recovery(),
		middleware.CORS(),
	)

	healthHandler := handlers.NewHealthHandler()

	v1 := router.Group("/api/v1")
	{
		public := v1.Group("")
		public.Use(container.RateLimitMiddleware.Handler("public"))
		{
			public.GET("/health", healthHandler.Health)

			auth := public.Group("/auth")
			{
				auth.POST("/register", container.AuthHandler.Register)
				auth.POST("/login", container.AuthHandler.Login)
				auth.POST("/refresh", container.AuthHandler.Refresh)
				auth.POST("/forgot-password", container.AuthHandler.ForgotPassword)
				auth.POST("/verify-email", container.AuthHandler.VerifyResetCode)
				auth.POST("/reset-password", container.AuthHandler.ResetPassword)
				auth.POST("/resend-verification-code", container.AuthHandler.ResendVerification)
			}
		}

		protected := v1.Group("")
		{
			protected.Use(
				container.AuthMiddleware.RequireAuth(),
				container.RateLimitMiddleware.Handler("protected"),
			)

			protected.POST("/auth/activate-account", container.AuthHandler.ActivateAccount)
			protected.POST("/auth/logout", container.AuthHandler.Logout)
			protected.POST("/auth/logout-all", container.AuthHandler.LogoutAllSessions)
		}

		verified := v1.Group("")
		{
			verified.Use(
				container.AuthMiddleware.RequireAuth(),
				middleware.RequireVerified(),
				container.RateLimitMiddleware.Handler("verified"),
			)

			user := verified.Group("/user")
			{
				user.GET("/me", container.UserHandler.GetAuthUser)
				user.PUT("/update-account", container.UserHandler.UpdateUserAccount)
				user.PUT("/update-password", container.UserHandler.UpdatePassword)
				user.PUT("/update-avatar", container.UserHandler.UpdateUserAvatar)

				user.GET("/audit-logs", container.UserHandler.GetAuditLogs)
				user.GET("/login-history", container.UserHandler.GetLoginHistory)
				user.GET("/sessions", container.UserHandler.GetUserSessions)
				user.GET("/sessions/current", container.UserHandler.GetCurrentSessions)
			}

			authorization := verified.Group("/authorization")
			{
				roles := authorization.Group("/roles")
				{
					roles.GET("/", container.AuthorizationHandler.GetRoles)
					roles.POST("/create", container.AuthorizationHandler.CreateRole)
					roles.PUT("/update/:id", container.AuthorizationHandler.UpdateRole)
					roles.DELETE("/delete/:id", container.AuthorizationHandler.DeleteRole)
					roles.GET("/get/:id", container.AuthorizationHandler.GetRole)
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

	return router
}
