package seed

import (
	"context"
	"time"

	appAuthorization "github.com/go-api/internal/application/authorization"
	"github.com/go-api/internal/bootstrap"
	backendUser "github.com/go-api/internal/domain/user"
	"github.com/go-api/internal/shared/config"
	"github.com/google/uuid"
)

func SeedUsers(ctx context.Context, container *bootstrap.Container, cfg *config.Config) error {
	user, err := container.UserRepository.FindByEmail(ctx, cfg.Admin.Email)
	if err != nil || user == nil {

		hash, err := container.PasswordService.Hash(cfg.Admin.Password)
		if err != nil {
			return err
		}

		now := time.Now()

		user = &backendUser.User{
			ID:         uuid.New(),
			FirstName:  cfg.Admin.FirstName,
			LastName:   cfg.Admin.LastName,
			Email:      cfg.Admin.Email,
			Password:   hash,
			IsVerified: true,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		if err := container.UserRepository.Create(ctx, user); err != nil {
			return err
		}
	}

	role, err := container.AuthorizationService.GetRoleByName(ctx, "Super Admin")
	if err != nil {
		return err
	}

	hasRole, err := container.AuthorizationService.UserHasRole(ctx, user.ID, role.ID)
	if err != nil {
		return err
	}

	if hasRole {
		return nil
	}

	return container.AuthorizationService.AssignRoleToUser(
		ctx,
		appAuthorization.AssignRoleRequest{
			UserID: user.ID,
			RoleID: role.ID,
		},
	)
}
