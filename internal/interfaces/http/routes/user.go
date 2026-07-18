package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api/internal/bootstrap"
	"github.com/go-api/internal/interfaces/http/middleware"
)

func RegisterUserRoutes(router *gin.RouterGroup, container *bootstrap.Container) {

	verified := router.Group("")
	verified.Use(
		container.AuthMiddleware.RequireAuth(),
		middleware.RequireVerified(),
		container.RateLimitMiddleware.Handler("verified"),
	)

	user := verified.Group("/user")

	user.GET("/me", container.UserHandler.GetAuthUser)
	user.PUT("/update-account", container.UserHandler.UpdateUserAccount)
	user.PUT("/update-password", container.UserHandler.UpdatePassword)
	user.PUT("/update-avatar", container.UserHandler.UpdateUserAvatar)

	user.GET("/audit-logs", container.UserHandler.GetAuditLogs)
	user.GET("/login-history", container.UserHandler.GetLoginHistory)
	user.GET("/sessions", container.UserHandler.GetUserSessions)
	user.GET("/sessions/current", container.UserHandler.GetCurrentSessions)
}
