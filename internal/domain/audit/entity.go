package audit

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID         uuid.UUID
	UserID     *uuid.UUID
	Action     string
	EntityType *string
	EntityID   *uuid.UUID
	IPAddress  net.IP
	UserAgent  string
	Metadata   map[string]any
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
