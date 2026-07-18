package seed

import (
	"context"

	"github.com/go-api/internal/bootstrap"
	"github.com/go-api/internal/shared/config"
)

func Run(ctx context.Context, container *bootstrap.Container, cfg *config.Config) error {

	if err := SeedPermissions(ctx, container.AuthorizationService); err != nil {
		return err
	}

	if err := SeedRoles(ctx, container.AuthorizationService); err != nil {
		return err
	}

	if err := SeedAssignments(ctx, container.AuthorizationService); err != nil {
		return err
	}

	if err := SeedUsers(ctx, container, cfg); err != nil {
		return err
	}

	return nil
}
