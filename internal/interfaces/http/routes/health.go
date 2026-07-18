package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/go-api/internal/interfaces/http/handlers"
)

func RegisterHealthRoutes(router *gin.RouterGroup) {

	health := handlers.NewHealthHandler()

	public := router.Group("")
	{
		public.GET("/health", health.Health)
	}
}
