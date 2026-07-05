package user

import (
	"context"

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

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}
