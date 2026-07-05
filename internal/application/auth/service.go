package auth

import (
	"context"
	"log"
	"net"
	"time"

	applicationToken "github.com/go-api/internal/application/token"
	backendSession "github.com/go-api/internal/domain/session"
	domainToken "github.com/go-api/internal/domain/token"
	backendUser "github.com/go-api/internal/domain/user"
	"github.com/go-api/internal/shared/config"
	"github.com/google/uuid"

	"github.com/go-api/internal/infrastructure/email"
	"github.com/go-api/internal/infrastructure/jwt"
)

type Service struct {
	cfg          *config.JWTConfig
	users        backendUser.Repository
	sessions     backendSession.Repository
	password     *PasswordService
	tokenService *applicationToken.Service
	emailService *email.Service
	jwtManager   *jwt.Manager
}

func NewService(
	cfg *config.Config,
	users backendUser.Repository,
	sessions backendSession.Repository,
	password *PasswordService,
	tokenService *applicationToken.Service,
	emailService *email.Service,
	jwtManager *jwt.Manager,
) *Service {

	return &Service{
		cfg:          &cfg.JWT,
		users:        users,
		sessions:     sessions,
		password:     password,
		tokenService: tokenService,
		emailService: emailService,
		jwtManager:   jwtManager,
	}
}

func (s *Service) AuthenticateUser(ctx context.Context, user *backendUser.User, userAgent, ipAddress string) (*LoginResponse, error) {
	refreshToken, err := s.tokenService.CreateRefreshToken()
	if err != nil {
		return nil, err
	}

	session := &backendSession.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken.TokenHash,
		UserAgent:    userAgent,
		IPAddress:    net.ParseIP(ipAddress),
		ExpiresAt:    refreshToken.ExpiresAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.sessions.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, session.ID)
	if err != nil {
		return nil, err
	}

	if !user.IsVerified {
		tokenType := domainToken.AccountActivation
		code, err := s.tokenService.CreateCodeToken(ctx, user.ID, tokenType)
		if err != nil {
			return nil, err
		}
		log.Println("Activation code sent successfully:", code)

		// subject, html, text := email.VerificationTemplate(user.FirstName, code)

		// err = s.emailService.Send(ctx, email.Message{
		// 	To:      user.Email,
		// 	Subject: subject,
		// 	HTML:    html,
		// 	Text:    text,
		// })

		// if err != nil {
		// 	fmt.Println("Email verification failed:", err)
		// }
	}

	return &LoginResponse{
		User:         NewUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    int(s.cfg.AccessTokenMinutes),
	}, nil
}

func (s *Service) Register(ctx context.Context, req RegisterRequest, userAgent, ipAddress string) (*LoginResponse, error) {

	existing, _ := s.users.FindByEmail(ctx, req.Email)

	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := s.password.Hash(req.Password)

	if err != nil {
		return nil, err
	}

	user := &backendUser.User{
		ID:         uuid.New(),
		Email:      req.Email,
		Password:   hash,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.users.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return s.AuthenticateUser(ctx, user, userAgent, ipAddress)
}

func (s *Service) Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (*LoginResponse, error) {

	user, err := s.users.FindByEmail(ctx, req.Email)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.password.Verify(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.AuthenticateUser(ctx, user, userAgent, ipAddress)
}

func (s *Service) ActivateAccount(ctx context.Context, userID uuid.UUID, req ActivateAccountRequest) error {

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return ErrUnauthorized
	}

	if user.IsVerified {
		return ErrEmailAlreadyVerified
	}

	hash := s.tokenService.Hash(req.Code)
	token, err := s.tokenService.FindByTokenAndTypeAndUser(
		ctx,
		hash,
		domainToken.AccountActivation,
		user.ID,
	)

	if err != nil {
		return ErrInvalidVerificationCode
	}

	if token.UserID != user.ID {
		return ErrInvalidVerificationCode
	}

	if token.UsedAt != nil {
		return ErrInvalidVerificationCode
	}

	if token.ExpiresAt.Before(time.Now()) {
		_ = s.tokenService.DeleteByUserAndType(ctx, user.ID, domainToken.AccountActivation)

		return ErrVerificationCodeExpired
	}

	err = s.users.Verify(ctx, user.ID)
	if err != nil {
		return err
	}

	err = s.tokenService.Consume(ctx, token.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return ErrUserNotFound
	}

	tokenType := domainToken.PasswordReset
	code, err := s.tokenService.CreateCodeToken(ctx, user.ID, tokenType)
	if err != nil {
		return err
	}

	log.Println("Reset password code sent successfully:", code)

	// subject, html, text := email.VerificationTemplate(user.FirstName, code)

	// err = s.emailService.Send(ctx, email.Message{
	// 	To:      user.Email,
	// 	Subject: subject,
	// 	HTML:    html,
	// 	Text:    text,
	// })

	// if err != nil {
	// 	fmt.Println("Email verification failed:", err)
	// }

	return nil
}

func (s *Service) ResendVerification(ctx context.Context, req ResendVerificationRequest) error {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return ErrUserNotFound
	}

	tokenType := domainToken.PasswordReset
	code, err := s.tokenService.CreateCodeToken(ctx, user.ID, tokenType)
	if err != nil {
		return err
	}

	log.Println("Resend verification code sent successfully:", code)

	// subject, html, text := email.VerificationTemplate(user.FirstName, code)

	// err = s.emailService.Send(ctx, email.Message{
	// 	To:      user.Email,
	// 	Subject: subject,
	// 	HTML:    html,
	// 	Text:    text,
	// })

	// if err != nil {
	// 	fmt.Println("Email verification failed:", err)
	// }

	return nil
}

func (s *Service) VerifyResetCode(ctx context.Context, req VerifyResetCodeRequest) (*VerifyResetCodeResponse, error) {

	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidVerificationCode
	}

	hash := s.tokenService.Hash(req.Code)
	token, err := s.tokenService.FindByTokenAndTypeAndUser(
		ctx,
		hash,
		domainToken.PasswordReset,
		user.ID,
	)

	if err != nil {
		return nil, ErrInvalidVerificationCode
	}

	if token.UserID != user.ID {
		return nil, ErrInvalidVerificationCode
	}

	if token.UsedAt != nil {
		return nil, ErrInvalidVerificationCode
	}

	if token.ExpiresAt.Before(time.Now()) {
		_ = s.tokenService.DeleteByUserAndType(ctx, user.ID, domainToken.PasswordReset)

		return nil, ErrVerificationCodeExpired
	}

	err = s.tokenService.Consume(ctx, token.ID)
	if err != nil {
		return nil, err
	}

	resetToken, err := s.tokenService.CreateResetToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &VerifyResetCodeResponse{
		ResetToken: resetToken,
	}, nil
}

func (s *Service) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	hash := s.tokenService.Hash(req.ResetToken)

	token, err := s.tokenService.FindByTokenAndType(
		ctx,
		hash,
		domainToken.PasswordResetGrant,
	)
	if err != nil {
		return ErrInvalidResetToken
	}

	if token.UsedAt != nil {
		return ErrResetTokenAlreadyUsed
	}

	if token.ExpiresAt.Before(time.Now()) {
		_ = s.tokenService.DeleteByUserAndType(ctx, token.UserID, domainToken.PasswordResetGrant)

		return ErrResetTokenExpired
	}

	user, err := s.users.FindByID(ctx, token.UserID)
	if err != nil {
		return ErrInvalidResetToken
	}

	password, err := s.password.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	if err := s.users.UpdatePassword(ctx, user.ID, password); err != nil {
		return err
	}

	if err := s.tokenService.Consume(ctx, token.ID); err != nil {
		return err
	}

	if err := s.sessions.RevokeAllSessions(
		ctx,
		user.ID,
		backendSession.RevokedByPasswordChange,
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (*LoginResponse, error) {

	hashed := s.tokenService.Hash(req.RefreshToken)

	session, err := s.sessions.FindByToken(ctx, hashed)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if session.RevokedAt != nil {
		return nil, ErrInvalidRefreshToken
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.sessions.RevokeRefreshToken(
			ctx,
			session.ID,
			backendSession.TokenExpired,
		)
		return nil, ErrSessionExpired
	}

	err = s.sessions.UpdateLastUsedAt(ctx, session.ID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.users.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, session.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.CreateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = s.sessions.UpdateRefreshToken(
		ctx,
		session.ID,
		refreshToken.TokenHash,
		refreshToken.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         NewUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    s.cfg.AccessTokenMinutes,
	}, nil
}

func (s *Service) Logout(ctx context.Context, sessionID uuid.UUID) error {

	err := s.sessions.RevokeRefreshToken(
		ctx,
		sessionID,
		backendSession.RevokedByLogout,
	)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	err = s.sessions.UpdateLastUsedAt(ctx, sessionID)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	return nil
}

func (s *Service) LogoutAllSessions(ctx context.Context, userID uuid.UUID) error {

	err := s.sessions.RevokeAllSessions(
		ctx,
		userID,
		backendSession.RevokedByLogoutAll,
	)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	return nil
}
