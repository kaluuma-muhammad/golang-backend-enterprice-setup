package auth

import (
	"context"
	"net"
	"time"

	backendSession "github.com/go-api/internal/domain/session"
	backendUser "github.com/go-api/internal/domain/user"
	"github.com/google/uuid"

	"github.com/go-api/internal/infrastructure/jwt"
)

const RefreshTTL = 7 * 24 * time.Hour

type Service struct {
	users        backendUser.Repository
	sessions     backendSession.Repository
	password     *PasswordService
	tokenService *TokenService
	jwtManager   *jwt.Manager
}

func NewService(
	users backendUser.Repository,
	sessions backendSession.Repository,
	password *PasswordService,
	tokenService *TokenService,
	jwtManager *jwt.Manager,
) *Service {

	return &Service{
		users:        users,
		sessions:     sessions,
		password:     password,
		tokenService: tokenService,
		jwtManager:   jwtManager,
	}
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

	refreshToken, err := s.tokenService.Generate()

	if err != nil {
		return nil, err
	}

	refreshHash := s.tokenService.Hash(refreshToken)

	session := &backendSession.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshHash,
		UserAgent:    userAgent,
		IPAddress:    net.ParseIP(ipAddress),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
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

	return &LoginResponse{
		User:         NewUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (*LoginResponse, error) {

	user, err := s.users.FindByEmail(ctx, req.Email)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.password.Verify(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	refreshToken, err := s.tokenService.Generate()

	if err != nil {
		return nil, err
	}

	refreshHash := s.tokenService.Hash(refreshToken)

	session := &backendSession.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshHash,
		UserAgent:    userAgent,
		IPAddress:    net.ParseIP(ipAddress),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
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

	return &LoginResponse{
		User:         NewUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	}, nil
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

	refreshToken, err := s.tokenService.Generate()
	if err != nil {
		return nil, err
	}

	refreshHash := s.tokenService.Hash(refreshToken)

	err = s.sessions.UpdateRefreshToken(
		ctx,
		session.ID,
		refreshHash,
		time.Now().Add(RefreshTTL),
	)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         NewUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(jwt.AccessTokenTTL.Seconds()),
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
