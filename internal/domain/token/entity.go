package token

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Type      Type
	Token     string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
