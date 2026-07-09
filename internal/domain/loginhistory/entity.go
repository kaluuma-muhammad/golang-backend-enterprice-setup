package loginhistory

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type LoginHistory struct {
	ID            uuid.UUID
	UserID        *uuid.UUID
	SessionID     *uuid.UUID
	Status        Status
	IPAddress     net.IP
	UserAgent     string
	DeviceName    string
	FailureReason *string
	LoginAt       time.Time
	LogoutAt      *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
