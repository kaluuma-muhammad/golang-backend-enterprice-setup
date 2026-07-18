package bootstrap

import (
	"log"

	"github.com/go-api/internal/application/auth"
	"github.com/go-api/internal/application/authorization"
	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/application/security"
	"github.com/go-api/internal/application/token"
	"github.com/go-api/internal/application/user"
	"github.com/go-api/internal/infrastructure/email"
	"github.com/go-api/internal/infrastructure/jwt"
	repositories "github.com/go-api/internal/infrastructure/postgres/repositories"
	"github.com/go-api/internal/infrastructure/storage"
	handlers "github.com/go-api/internal/interfaces/http/handlers"
	authorizationHandlers "github.com/go-api/internal/interfaces/http/handlers/authorization"
	"github.com/go-api/internal/interfaces/http/middleware"
	"github.com/go-api/internal/shared/config"
	"github.com/go-api/internal/shared/ratelimiter"
)

func registerAuthDependencies(container *Container, cfg *config.Config) {
	config := ratelimiter.DefaultConfig()
	store := ratelimiter.NewMemoryStore()

	rateLimiter := ratelimiter.NewService(config, store)
	rateLimiter.StartCleanup()

	container.RateLimiter = rateLimiter
	container.RateLimitMiddleware = middleware.NewRateLimitMiddleware(rateLimiter)

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
	roleRepo := repositories.NewRoleRepository(container.DB)
	permissionRepo := repositories.NewPermissionRepository(container.DB)
	assignmentRepo := repositories.NewAssignmentRepository(container.DB)

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

	authorizationService := authorization.NewService(
		roleRepo,
		permissionRepo,
		assignmentRepo,
		userRepo,
	)

	authService := auth.NewService(
		cfg,
		userRepo,
		sessionRepo,
		passwordService,
		tokenService,
		emailService,
		authorizationService,
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

	authorizationMiddleware := middleware.NewAuthorizationMiddleware(authorizationService)

	container.Authenticator = authenticator
	container.SecurityService = securityService
	container.AuthorizationService = authorizationService
	container.AuthorizationMiddleware = authorizationMiddleware

	container.PasswordService = passwordService
	container.UserRepository = userRepo
	container.RoleRepository = roleRepo
	container.PermissionRepository = permissionRepo
	container.AssignmentRepository = assignmentRepo

	// HTTP Handlers
	container.AuthMiddleware = middleware.NewAuthMiddleware(authenticator)
	container.AuthHandler = handlers.NewAuthHandler(authService, container.Validator)
	container.UserHandler = handlers.NewUserHandler(userService, securityService, container.Validator)
	container.AuthorizationHandler = authorizationHandlers.NewAuthorizationHandler(
		authorizationService,
		container.Validator,
	)
}
