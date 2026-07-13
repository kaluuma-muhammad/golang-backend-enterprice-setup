package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID
	Email               string
	Password            string
	FirstName           string
	LastName            string
	Phone               *string
	ImageURL            *string
	IsVerified          bool
	FailedLoginAttempts int
	LockedUntil         *time.Time
	LastLoginAt         *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
