package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	appAuth "github.com/go-api/internal/application/auth"
	"github.com/go-api/internal/interfaces/http/response"
)

const (
	ContextUserKey    = "current_user"
	ContextSessionKey = "current_session"
)

type AuthMiddleware struct {
	authenticator *appAuth.Authenticator
}

func NewAuthMiddleware(
	authenticator *appAuth.Authenticator,
) *AuthMiddleware {

	return &AuthMiddleware{
		authenticator: authenticator,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" {
			response.Error(
				c,
				response.StatusCode(appAuth.ErrMissingAuthorizationHeader),
				appAuth.ErrMissingAuthorizationHeader.Error(),
			)
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(
				c,
				response.StatusCode(appAuth.ErrInvalidAuthorizationHeader),
				appAuth.ErrInvalidAuthorizationHeader.Error(),
			)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		authenticated, err := m.authenticator.Authenticate(c.Request.Context(), token)

		if err != nil {
			response.Error(
				c,
				response.StatusCode(err),
				err.Error(),
			)
			c.Abort()
			return
		}

		c.Set(ContextUserKey, authenticated.User)
		c.Set(ContextSessionKey, authenticated.Session)

		c.Next()
	}
}

func RequireVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)

		if user == nil {
			response.Error(c, http.StatusUnauthorized, appAuth.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		if !user.IsVerified {
			response.Error(c, http.StatusForbidden, appAuth.ErrEmailNotVerified.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}

func IsVerified(c *gin.Context) bool {

	user := CurrentUser(c)

	if user == nil {
		return false
	}

	return user.IsVerified
}
