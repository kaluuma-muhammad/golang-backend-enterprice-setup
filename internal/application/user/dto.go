package user

import "io"

type UpdateUserAccountRequest struct {
	Email     string  `json:"email" binding:"required"`
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Phone     *string `json:"phone" binding:"omitempty"`
}

type UploadAvatarRequest struct {
	File        io.Reader `form:"file" binding:"required"`
	Filename    string    `form:"filename" binding:"required"`
	ContentType string    `form:"content_type" binding:"required"`
}

type UpdateUserPasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
