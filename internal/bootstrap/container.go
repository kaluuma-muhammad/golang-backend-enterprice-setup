package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-api/internal/application/auth"
	handlers "github.com/go-api/internal/interfaces/http/handlers"
	"github.com/go-api/internal/interfaces/http/middleware"
	"github.com/go-api/internal/shared/config"
	"github.com/go-api/internal/shared/validator"
)

type Container struct {
	DB             *pgxpool.Pool
	Validator      *validator.Validator
	AuthHandler    *handlers.AuthHandler
	Authenticator  *auth.Authenticator
	AuthMiddleware *middleware.AuthMiddleware
}

func NewContainer(dbPool *pgxpool.Pool, cfg *config.Config) *Container {

	container := &Container{
		DB:        dbPool,
		Validator: validator.New(),
	}

	registerAuthDependencies(container, cfg)

	return container
}
