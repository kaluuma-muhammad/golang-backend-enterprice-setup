package token

import (
	"time"

	domainToken "github.com/go-api/internal/domain/token"
	"github.com/google/uuid"
)

type CreateTokenRequest struct {
	UserID uuid.UUID
	Type   domainToken.Type
	TTL    time.Duration
}

type VerifyTokenRequest struct {
	Token string
	Type  domainToken.Type
}
