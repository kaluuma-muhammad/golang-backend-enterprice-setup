package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-api/internal/interfaces/http/response"
)

func Recovery() gin.HandlerFunc {

	return gin.CustomRecovery(
		func(c *gin.Context, recovered any) {

			response.Error(c, http.StatusInternalServerError, "internal server error")
		},
	)
}
