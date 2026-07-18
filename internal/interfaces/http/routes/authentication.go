package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-api/internal/bootstrap"
)

func RegisterAuthenticationRoutes(router *gin.RouterGroup, container *bootstrap.Container) {

	public := router.Group("")
	public.Use(container.RateLimitMiddleware.Handler("public"))
	{
		auth := public.Group("/auth")

		auth.POST("/register", container.AuthHandler.Register)
		auth.POST("/login", container.AuthHandler.Login)
		auth.POST("/refresh", container.AuthHandler.Refresh)
		auth.POST("/forgot-password", container.AuthHandler.ForgotPassword)
		auth.POST("/verify-email", container.AuthHandler.VerifyResetCode)
		auth.POST("/reset-password", container.AuthHandler.ResetPassword)
		auth.POST("/resend-verification-code", container.AuthHandler.ResendVerification)
	}

	protected := router.Group("")
	protected.Use(
		container.AuthMiddleware.RequireAuth(),
		container.RateLimitMiddleware.Handler("protected"),
	)
	{
		auth := protected.Group("/auth")

		auth.POST("/activate-account", container.AuthHandler.ActivateAccount)
		auth.POST("/logout", container.AuthHandler.Logout)
		auth.POST("/logout-all", container.AuthHandler.LogoutAllSessions)
	}
}
