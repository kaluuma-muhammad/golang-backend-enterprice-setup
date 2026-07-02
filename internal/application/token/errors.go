package token

import "errors"

var (
	ErrTokenNotFound         = errors.New("token not found")
	ErrTokenExpired          = errors.New("token expired")
	ErrTokenAlreadyUsed      = errors.New("token already used")
	ErrInvalidToken          = errors.New("invalid token")
	ErrInvalidTokenType      = errors.New("invalid token type")
	ErrFailedToGenerateToken = errors.New("failed to generate token")
)
