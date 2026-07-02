package token

import "time"

type RefreshTokenResult struct {
	Token     string
	TokenHash string
	ExpiresAt time.Time
}
