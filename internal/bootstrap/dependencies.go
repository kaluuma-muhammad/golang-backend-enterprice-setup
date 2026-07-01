package bootstrap

import (
	"github.com/go-api/internal/application/auth"
	"github.com/go-api/internal/infrastructure/jwt"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
	repositories "github.com/go-api/internal/infrastructure/postgres/repositories"
	handlers "github.com/go-api/internal/interfaces/http/handlers"
	"github.com/go-api/internal/interfaces/http/middleware"
	"github.com/go-api/internal/shared/config"
)

func registerAuthDependencies(container *Container, cfg *config.Config) {

	// Database
	queries := db.New(container.DB)
	jwtManager := jwt.NewManager(cfg.JWT.Secret)

	// Repositories
	userRepo := repositories.NewUserRepository(queries)
	sessionRepo := repositories.NewSessionRepository(queries)

	// Application Services
	tokenService := auth.NewTokenService()
	passwordService := auth.NewPasswordService()

	authService := auth.NewService(
		userRepo,
		sessionRepo,
		passwordService,
		tokenService,
		jwtManager,
	)

	authenticator := auth.NewAuthenticator(
		userRepo,
		sessionRepo,
		jwtManager,
	)

	// HTTP Handlers
	container.AuthHandler = handlers.NewAuthHandler(authService, container.Validator)
	container.Authenticator = authenticator
	container.AuthMiddleware = middleware.NewAuthMiddleware(authenticator)

}
