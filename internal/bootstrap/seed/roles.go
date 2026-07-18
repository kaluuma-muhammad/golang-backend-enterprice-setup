package seed

import (
	"context"
	"errors"

	appAuthorization "github.com/go-api/internal/application/authorization"
)

var roles = []appAuthorization.CreateRoleRequest{
	{
		Name:        "Super Admin",
		Description: "Has unrestricted access to the entire system.",
		IsSystem:    true,
	},
	{
		Name:        "Admin",
		Description: "Manages users and application resources.",
		IsSystem:    true,
	},
	{
		Name:        "User",
		Description: "Default role assigned to authenticated users.",
		IsSystem:    true,
	},
}

func SeedRoles(ctx context.Context, service *appAuthorization.Service) error {
	for _, role := range roles {
		_, err := service.GetRoleByName(ctx, role.Name)
		if err == nil {
			continue
		}

		_, err = service.CreateRole(ctx, role)
		if err != nil &&
			!errors.Is(err, appAuthorization.ErrRoleAlreadyExists) {
			return err
		}
	}

	return nil
}
