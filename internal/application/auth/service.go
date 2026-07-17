package auth

import (
	"strings"

	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/application/security"
	applicationToken "github.com/go-api/internal/application/token"
	backendSession "github.com/go-api/internal/domain/session"
	backendUser "github.com/go-api/internal/domain/user"
	"github.com/go-api/internal/shared/config"

	"github.com/go-api/internal/infrastructure/email"
	"github.com/go-api/internal/infrastructure/jwt"
)

type Service struct {
	cfg          *config.JWTConfig
	users        backendUser.Repository
	sessions     backendSession.Repository
	password     *security.PasswordService
	tokenService *applicationToken.Service
	emailService *email.Service
	jwtManager   *jwt.Manager
	security     *security.Service
	tx           common.TransactionManager
	baseURL      string
}

func NewService(
	cfg *config.Config,
	users backendUser.Repository,
	sessions backendSession.Repository,
	password *security.PasswordService,
	tokenService *applicationToken.Service,
	emailService *email.Service,
	jwtManager *jwt.Manager,
	security *security.Service,
	tx common.TransactionManager,
	baseURL string,
) *Service {

	return &Service{
		cfg:          &cfg.JWT,
		users:        users,
		sessions:     sessions,
		password:     password,
		tokenService: tokenService,
		emailService: emailService,
		jwtManager:   jwtManager,
		security:     security,
		tx:           tx,
		baseURL:      baseURL,
	}
}

func parseDeviceName(userAgent string) string {
	if userAgent == "" {
		return "Unknown Device"
	}

	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "android"):
		return "Android Device"
	case strings.Contains(ua, "iphone"):
		return "iPhone"
	case strings.Contains(ua, "ipad"):
		return "iPad"
	case strings.Contains(ua, "windows"):
		return "Windows PC"
	case strings.Contains(ua, "mac"):
		return "Mac"
	case strings.Contains(ua, "linux") || strings.Contains(ua, "x11"):
		return "Linux PC"
	default:
		return "Unknown Device"
	}
}
