package entities

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID
	Resource    string
	Action      string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
