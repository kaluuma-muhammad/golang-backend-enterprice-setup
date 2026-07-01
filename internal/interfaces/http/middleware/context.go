package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	backendSession "github.com/go-api/internal/domain/session"
	backendUser "github.com/go-api/internal/domain/user"
)

func CurrentUser(c *gin.Context) *backendUser.User {
	value, exists := c.Get(ContextUserKey)
	if !exists {
		return nil
	}

	user, ok := value.(*backendUser.User)
	if !ok {
		return nil
	}

	return user
}

func CurrentSession(c *gin.Context) *backendSession.Session {
	value, exists := c.Get(ContextSessionKey)
	if !exists {
		return nil
	}

	session, ok := value.(*backendSession.Session)
	if !ok {
		return nil
	}

	return session
}

func UserID(c *gin.Context) uuid.UUID {
	user := CurrentUser(c)
	if user == nil {
		return uuid.Nil
	}

	return user.ID
}

func SessionID(c *gin.Context) uuid.UUID {
	session := CurrentSession(c)
	if session == nil {
		return uuid.Nil
	}

	return session.ID
}
