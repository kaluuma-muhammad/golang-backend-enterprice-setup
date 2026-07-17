package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appAuth "github.com/go-api/internal/application/auth"
	appAuthorization "github.com/go-api/internal/application/authorization"
	"github.com/go-api/internal/interfaces/http/response"
)

type AuthorizationMiddleware struct {
	authorization *appAuthorization.Service
}

func NewAuthorizationMiddleware(
	authorization *appAuthorization.Service,
) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		authorization: authorization,
	}
}

// Ensure that the authenticated user has the specified permission.
func (m *AuthorizationMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil {
			response.Error(c, http.StatusUnauthorized, appAuth.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		if err := m.authorization.Authorize(c.Request.Context(), user.ID, permission); err != nil {
			response.Error(c, response.StatusCode(err), err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}

// Ensure that the authenticated user has at least one of the provided permissions.
func (m *AuthorizationMiddleware) RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil {
			response.Error(c, http.StatusUnauthorized, appAuth.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		if err := m.authorization.AuthorizeAny(c.Request.Context(), user.ID, permissions...); err != nil {
			response.Error(c, response.StatusCode(err), err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}

// Ensure that the authenticated user has every permission supplied.
func (m *AuthorizationMiddleware) RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := CurrentUser(c)
		if user == nil {
			response.Error(c, http.StatusUnauthorized, appAuth.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		if err := m.authorization.AuthorizeAll(c.Request.Context(), user.ID, permissions...); err != nil {
			response.Error(c, response.StatusCode(err), err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
