package bootstrap

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-api/internal/application/auth"
	"github.com/go-api/internal/application/authorization"
	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/application/security"
	authorizationRepo "github.com/go-api/internal/domain/authorization/repositories"
	userRepo "github.com/go-api/internal/domain/user"
	redisinfra "github.com/go-api/internal/infrastructure/redis"
	handlers "github.com/go-api/internal/interfaces/http/handlers"
	authorizationHandler "github.com/go-api/internal/interfaces/http/handlers/authorization"
	"github.com/go-api/internal/interfaces/http/middleware"
	"github.com/go-api/internal/shared/config"
	"github.com/go-api/internal/shared/ratelimiter"
	"github.com/go-api/internal/shared/validator"
)

type Container struct {
	DB                      *pgxpool.Pool
	Validator               *validator.Validator
	Redis                   *redisinfra.Client
	Cache                   common.Cache
	AuthHandler             *handlers.AuthHandler
	UserHandler             *handlers.UserHandler
	AuthorizationHandler    *authorizationHandler.AuthorizationHandler
	Authenticator           *auth.Authenticator
	RateLimiter             ratelimiter.RateLimiter
	AuthMiddleware          *middleware.AuthMiddleware
	AuthorizationMiddleware *middleware.AuthorizationMiddleware
	RateLimitMiddleware     *middleware.RateLimitMiddleware
	SecurityService         *security.Service
	PasswordService         *security.PasswordService
	AuthorizationService    *authorization.Service
	UserRepository          userRepo.Repository
	RoleRepository          authorizationRepo.RoleRepository
	PermissionRepository    authorizationRepo.PermissionRepository
	AssignmentRepository    authorizationRepo.AssignmentRepository
}

func NewContainer(dbPool *pgxpool.Pool, cfg *config.Config) (*Container, error) {

	container := &Container{
		DB:        dbPool,
		Validator: validator.New(),
	}

	err := registerAuthDependencies(container, cfg)
	if err != nil {
		return nil, fmt.Errorf("register auth dependencies: %w", err)
	}

	return container, nil
}

func (a *App) Close() error {
	if a.Container != nil {
		return a.Container.Close()
	}

	return nil
}
