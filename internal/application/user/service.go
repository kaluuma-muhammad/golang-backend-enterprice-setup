package user

import (
	"context"
	"errors"

	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/application/security"
	"github.com/go-api/internal/application/token"
	domainToken "github.com/go-api/internal/domain/token"
	backendUser "github.com/go-api/internal/domain/user"
	"github.com/go-api/internal/infrastructure/email"
	"github.com/go-api/internal/infrastructure/storage"
	"github.com/google/uuid"
)

type Service struct {
	users        backendUser.Repository
	tx           common.TransactionManager
	password     *security.PasswordService
	emailService *email.Service
	tokenService *token.Service
	storage      storage.Storage
	baseURL      string
}

func NewService(
	users backendUser.Repository,
	tx common.TransactionManager,
	password *security.PasswordService,
	emailService *email.Service,
	tokenService *token.Service,
	storage storage.Storage,
	baseURL string,
) *Service {
	return &Service{
		users:        users,
		tx:           tx,
		password:     password,
		emailService: emailService,
		tokenService: tokenService,
		storage:      storage,
		baseURL:      baseURL,
	}
}

func (s *Service) GetAuthUser(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return NewUserResponse(user, s.baseURL), nil
}

func (s *Service) UpdateUserAccount(ctx context.Context, userID uuid.UUID, req UpdateUserAccountRequest) (*UserResponse, error) {
	var response *UserResponse

	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		user, err := s.users.FindByID(txCtx, userID)
		if err != nil {
			return err
		}

		if user.Email != req.Email {
			user.IsVerified = false

			tokenType := domainToken.EmailChange
			code, err := s.tokenService.CreateCodeToken(txCtx, user.ID, tokenType)
			if err != nil {
				return err
			}

			go s.emailService.SendVerificationEmail(context.Background(), user, code)
		}

		user.Email = req.Email
		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.Phone = req.Phone

		if err := s.users.UpdateUserAccount(txCtx, user); err != nil {
			return err
		}

		response = NewUserResponse(user, s.baseURL)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) UpdatePassword(ctx context.Context, userID uuid.UUID, req UpdateUserPasswordRequest) (*UserResponse, error) {
	var response *UserResponse

	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		user, err := s.users.FindByID(txCtx, userID)
		if err != nil {
			return err
		}

		if !s.password.Verify(req.CurrentPassword, user.Password) {
			return errors.New("Incorrect password")
		}

		if req.NewPassword != req.ConfirmPassword {
			return errors.New("New passwords do not match")
		}

		if req.CurrentPassword == req.NewPassword {
			return errors.New("New password cannot be the same as old password")
		}

		hashed, err := s.password.Hash(req.NewPassword)
		if err != nil {
			return err
		}

		if err := s.users.UpdatePassword(txCtx, user.ID, hashed); err != nil {
			return err
		}

		user.Password = hashed
		response = NewUserResponse(user, s.baseURL)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) UpdateUserAvatar(ctx context.Context, userID uuid.UUID, req UploadAvatarRequest) (*UserResponse, error) {
	var response *UserResponse

	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		user, err := s.users.FindByID(txCtx, userID)
		if err != nil {
			return err
		}

		uploaded, err := s.storage.Upload(storage.UploadInput{
			File:        req.File,
			Filename:    req.Filename,
			Folder:      "users",
			ContentType: req.ContentType,
		})
		if err != nil {
			return err
		}

		if user.ImageURL != nil && *user.ImageURL != "users/default-avatar.jpg" {
			if err := s.storage.Delete(*user.ImageURL); err != nil {
				return err
			}
		}

		user.ImageURL = &uploaded.Path

		if err := s.users.UpdateUserAvatar(txCtx, user); err != nil {
			return err
		}

		response = NewUserResponse(user, s.baseURL)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
