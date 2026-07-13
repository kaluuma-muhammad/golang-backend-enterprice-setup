package bootstrap

import (
	"log"

	"github.com/go-api/internal/application/auth"
	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/application/security"
	"github.com/go-api/internal/application/token"
	"github.com/go-api/internal/application/user"
	"github.com/go-api/internal/infrastructure/email"
	"github.com/go-api/internal/infrastructure/jwt"
	repositories "github.com/go-api/internal/infrastructure/postgres/repositories"
	"github.com/go-api/internal/infrastructure/storage"
	handlers "github.com/go-api/internal/interfaces/http/handlers"
	"github.com/go-api/internal/interfaces/http/middleware"
	"github.com/go-api/internal/shared/config"
)

func registerAuthDependencies(container *Container, cfg *config.Config) {
	jwtManager := jwt.NewManager(cfg.JWT)
	generator := token.NewGenerator()

	renderer, err := email.NewRenderer()
	if err != nil {
		log.Fatal(err)
	}

	provider, err := email.NewSMTP(
		cfg.Email.Host,
		cfg.Email.Port,
		cfg.Email.Username,
		cfg.Email.Password,
		cfg.Email.FromAddress,
		cfg.Email.FromName,
		cfg.Email.Encryption,
	)
	if err != nil {
		log.Fatal(err)
	}

	tx := common.NewTransactionManager(container.DB)

	// Repositories
	userRepo := repositories.NewUserRepository(container.DB)
	sessionRepo := repositories.NewSessionRepository(container.DB)
	tokenRepository := repositories.NewTokenRepository(container.DB)
	auditRepo := repositories.NewAuditRepository(container.DB)
	loginHistoryRepo := repositories.NewLoginHistoryRepository(container.DB)

	// Application Services
	emailService := email.New(cfg, provider, renderer)
	tokenService := token.NewService(tokenRepository, generator, &cfg.JWT)
	passwordService := security.NewPasswordService()
	authenticator := auth.NewAuthenticator(userRepo, sessionRepo, jwtManager)

	storageService := storage.NewLocalStorage(
		"./storage/images",
		cfg.App.BaseURL+"/storage/images",
	)

	securityService := security.NewService(
		userRepo,
		sessionRepo,
		auditRepo,
		loginHistoryRepo,
	)

	authService := auth.NewService(
		cfg,
		userRepo,
		sessionRepo,
		passwordService,
		tokenService,
		emailService,
		jwtManager,
		securityService,
		tx,
		cfg.App.BaseURL,
	)

	userService := user.NewService(
		userRepo,
		tx,
		passwordService,
		emailService,
		tokenService,
		storageService,
		cfg.App.BaseURL,
	)

	// HTTP Handlers
	container.AuthHandler = handlers.NewAuthHandler(authService, container.Validator)
	container.Authenticator = authenticator
	container.AuthMiddleware = middleware.NewAuthMiddleware(authenticator)
	container.SecurityService = securityService
	container.UserHandler = handlers.NewUserHandler(userService, securityService, container.Validator)

}
