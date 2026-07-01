package auth

import (
	"context"
	"time"

	backendSession "github.com/go-api/internal/domain/session"
	backendUser "github.com/go-api/internal/domain/user"
	"github.com/go-api/internal/infrastructure/jwt"
)

type AuthenticatedUser struct {
	User    *backendUser.User
	Session *backendSession.Session
}

type Authenticator struct {
	users      backendUser.Repository
	sessions   backendSession.Repository
	jwtManager *jwt.Manager
}

func NewAuthenticator(
	users backendUser.Repository,
	sessions backendSession.Repository,
	jwtManager *jwt.Manager,
) *Authenticator {

	return &Authenticator{
		users:      users,
		sessions:   sessions,
		jwtManager: jwtManager,
	}
}

func (a *Authenticator) Authenticate(
	ctx context.Context,
	accessToken string,
) (*AuthenticatedUser, error) {

	claims, err := a.jwtManager.Validate(accessToken)
	if err != nil {
		return nil, ErrUnauthorized
	}

	session, err := a.sessions.FindByID(ctx, claims.SessionID)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.RevokedAt != nil {
		return nil, ErrSessionRevoked
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, ErrSessionExpired
	}

	user, err := a.users.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, ErrUnauthorized
	}

	return &AuthenticatedUser{
		User:    user,
		Session: session,
	}, nil
}
