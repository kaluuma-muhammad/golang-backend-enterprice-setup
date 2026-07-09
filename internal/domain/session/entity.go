package session

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	RefreshToken  string
	UserAgent     string
	IPAddress     net.IP
	DeviceID      string
	Platform      string
	Browser       string
	LastSeenAt    time.Time
	LastUsedAt    *time.Time
	ExpiresAt     time.Time
	RevokedAt     *time.Time
	RevokedReason *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
