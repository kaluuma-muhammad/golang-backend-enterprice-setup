package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-api/internal/application/auth"
	"github.com/go-api/internal/application/authorization"
	"github.com/go-api/internal/application/security"
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
	AuthHandler             *handlers.AuthHandler
	UserHandler             *handlers.UserHandler
	AuthorizationHandler    *authorizationHandler.AuthorizationHandler
	Authenticator           *auth.Authenticator
	RateLimiter             *ratelimiter.Service
	AuthMiddleware          *middleware.AuthMiddleware
	AuthorizationMiddleware *middleware.AuthorizationMiddleware
	RateLimitMiddleware     *middleware.RateLimitMiddleware
	SecurityService         *security.Service
	AuthorizationService    *authorization.Service
}

func NewContainer(dbPool *pgxpool.Pool, cfg *config.Config) *Container {

	container := &Container{
		DB:        dbPool,
		Validator: validator.New(),
	}

	registerAuthDependencies(container, cfg)

	return container
}
