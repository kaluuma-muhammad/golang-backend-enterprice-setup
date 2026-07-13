package user

import (
	"context"

	"github.com/google/uuid"
)

type UserServiceContract interface {
	GetAuthUser(
		ctx context.Context,
		userID uuid.UUID,
	) (*UserResponse, error)

	UpdateUserAccount(
		ctx context.Context,
		userID uuid.UUID,
		req UpdateUserAccountRequest,
	) (*UserResponse, error)

	UpdatePassword(
		ctx context.Context,
		userID uuid.UUID,
		req UpdateUserPasswordRequest,
	) (*UserResponse, error)

	UpdateUserAvatar(
		ctx context.Context,
		userID uuid.UUID,
		req UploadAvatarRequest,
	) (*UserResponse, error)
}
