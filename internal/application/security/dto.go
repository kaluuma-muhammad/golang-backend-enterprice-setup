package security

import (
	"time"

	"github.com/google/uuid"
)

type AuditLogResponse struct {
	ID         uuid.UUID `json:"id"`
	Action     string    `json:"action"`
	EntityType *string   `json:"entity_type"`
	CreatedAt  time.Time `json:"created_at"`
}

type SessionResponse struct {
	ID         uuid.UUID  `json:"id"`
	Platform   string     `json:"platform"`
	Browser    string     `json:"browser"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
}

type LoginHistoryResponse struct {
	ID         uuid.UUID  `json:"id"`
	Status     string     `json:"status"`
	IPAddress  string     `json:"ip_address"`
	DeviceName string     `json:"device_name"`
	LoginAt    time.Time  `json:"login_at"`
	LogoutAt   *time.Time `json:"logout_at"`
}
