package auth

import (
	"errors"
)

var (
	// Login/Register
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")

	// Refresh Token
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	// sessions
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionRevoked  = errors.New("session has been revoked")
	ErrSessionExpired  = errors.New("session has expired")

	// Authentication
	ErrUnauthorized               = errors.New("unauthorized")
	ErrMissingAuthorizationHeader = errors.New("missing authorization header")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")

	// Account
	ErrEmailNotVerified = errors.New("email not verified")
)
