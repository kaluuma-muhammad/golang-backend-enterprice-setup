package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		user *User,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*User, error)

	FindByEmail(
		ctx context.Context,
		email string,
	) (*User, error)

	Update(
		ctx context.Context,
		user *User,
	) error

	Verify(
		ctx context.Context,
		id uuid.UUID,
	) error

	UpdatePassword(
		ctx context.Context,
		id uuid.UUID,
		password string,
	) error

	// Security

	IncrementFailedLoginAttempts(
		ctx context.Context,
		id uuid.UUID,
	) error

	ResetFailedLoginAttempts(
		ctx context.Context,
		id uuid.UUID,
	) error

	LockUserAccount(
		ctx context.Context,
		id uuid.UUID,
		until time.Time,
	) error

	UpdateLastLogin(
		ctx context.Context,
		id uuid.UUID,
	) error

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}
