package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Email    string
	Password string

	FirstName string
	LastName  string

	IsVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
