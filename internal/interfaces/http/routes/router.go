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
		}

	}

	return router
}
