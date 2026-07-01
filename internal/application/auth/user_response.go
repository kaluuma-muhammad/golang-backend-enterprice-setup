package auth

import (
	"github.com/go-api/internal/domain/user"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	IsVerified bool      `json:"is_verified"`
}

func NewUserResponse(u *user.User) UserResponse {
	return UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		IsVerified: u.IsVerified,
	}
}
