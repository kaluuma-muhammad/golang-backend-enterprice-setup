package user

import (
	"github.com/go-api/internal/domain/user"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Phone      *string   `json:"phone"`
	ImageURL   *string   `json:"image_url"`
	IsVerified bool      `json:"is_verified"`
}

func NewUserResponse(u *user.User, baseURL string) *UserResponse {
	image := baseURL + "/storage/images/users/default-avatar.jpg"

	if u.ImageURL != nil && *u.ImageURL != "" {
		image = baseURL + "/storage/images/" + *u.ImageURL
	}

	return &UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Phone:      u.Phone,
		ImageURL:   &image,
		IsVerified: u.IsVerified,
	}
}
