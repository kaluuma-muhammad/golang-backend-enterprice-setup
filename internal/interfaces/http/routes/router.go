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
		{
			public.GET("/health", healthHandler.Health)

			auth := public.Group("/auth")
			{
				auth.POST("/register", container.AuthHandler.Register)
				auth.POST("/login", container.AuthHandler.Login)
				auth.POST("/refresh", container.AuthHandler.Refresh)
			}
		}

		protected := v1.Group("")
		{
			protected.Use(container.AuthMiddleware.RequireAuth())

			// protected.POST("/auth/activate-account", container.AuthHandler.ActivateAccount)
			// protected.POST("/auth/verify-email", container.AuthHandler.VerifyEmail)
			// protected.POST("/auth/resend-verification", container.AuthHandler.ResendVerificationEmail)
			// protected.POST("/auth/reset-password", container.AuthHandler.ResetPassword)
			protected.POST("/auth/logout", container.AuthHandler.Logout)
			protected.POST("/auth/logout-all", container.AuthHandler.LogoutAllSessions)
		}

		verified := v1.Group("")
		{
			verified.Use(container.AuthMiddleware.RequireAuth(), middleware.RequireVerified())

			// account := verified.Group("/account")
			// {
			// 	account.GET("/me", container.AccountHandler.GetProfile)
			// 	account.PATCH("/me", container.AccountHandler.UpdateProfile)
			// }
		}

	}

	return router
}
