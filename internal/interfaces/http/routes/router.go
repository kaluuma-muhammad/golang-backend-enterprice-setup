package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/go-api/internal/bootstrap"
	"github.com/go-api/internal/interfaces/http/middleware"
)

func SetupRouter(logger *zap.Logger, container *bootstrap.Container) *gin.Engine {

	router := gin.New()

	router.Use(middleware.RequestID(), middleware.Logger(logger), middleware.Recovery(), middleware.CORS())

	v1 := router.Group("/api/v1")

	RegisterHealthRoutes(v1)
	RegisterAuthenticationRoutes(v1, container)
	RegisterUserRoutes(v1, container)
	RegisterAuthorizationRoutes(v1, container)

	return router
}
