package jwt

import "errors"

var (
	ErrTokenInvalid = errors.New("invalid token")

	ErrTokenInvalidClaims = errors.New("invalid token claims")
)
