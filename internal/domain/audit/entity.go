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
	Metadata   []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
