package auth

import (
	"context"
	"net"
	"time"

	userResponse "github.com/go-api/internal/application/user"
	backendSession "github.com/go-api/internal/domain/session"
	domainToken "github.com/go-api/internal/domain/token"
	backendUser "github.com/go-api/internal/domain/user"
	"github.com/google/uuid"
)

func (s *Service) AuthenticateUser(ctx context.Context, user *backendUser.User, userAgent, ipAddress string) (*LoginResponse, error) {
	deviceInfo := ParseDeviceInfo(userAgent, ipAddress)
	now := time.Now()

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
		DeviceID:     deviceInfo.DeviceID,
		Platform:     deviceInfo.Platform,
		Browser:      deviceInfo.Browser,
		LastUsedAt:   &now,
		ExpiresAt:    refreshToken.ExpiresAt,
		CreatedAt:    now,
		UpdatedAt:    now,
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

		go s.emailService.SendVerificationEmail(context.Background(), user, code)
	}

	err = s.security.RecordSuccessfulLogin(
		ctx,
		user.ID,
		session.ID,
		ipAddress,
		userAgent,
		deviceInfo.DeviceName,
	)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         userResponse.NewUserResponse(user, s.baseURL),
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    int(s.cfg.AccessTokenMinutes),
	}, nil
}

func (s *Service) Register(ctx context.Context, req RegisterRequest, userAgent, ipAddress string) (*LoginResponse, error) {

	var response *LoginResponse

	err := s.tx.Execute(ctx, func(txCtx context.Context) error {

		existing, _ := s.users.FindByEmail(txCtx, req.Email)

		if existing != nil {
			return ErrEmailAlreadyExists
		}

		hash, err := s.password.Hash(req.Password)

		if err != nil {
			return err
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

		err = s.users.Create(txCtx, user)

		if err != nil {
			return err
		}

		response, err = s.AuthenticateUser(txCtx, user, userAgent, ipAddress)
		return err
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (*LoginResponse, error) {
	var response *LoginResponse

	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		user, err := s.users.FindByEmail(txCtx, req.Email)
		if err != nil {
			return ErrInvalidCredentials
		}

		if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
			return ErrAccountLocked
		}

		if !s.password.Verify(req.Password, user.Password) {
			if err := s.security.IncrementFailedLoginAttempt(txCtx, user.ID); err != nil {
				return err
			}
			if err := s.security.RecordFailedLogin(txCtx, &user.ID, ipAddress, userAgent, "Invalid Credentials"); err != nil {
				return err
			}

			if user.FailedLoginAttempts >= 5 {
				if err := s.security.LockAccount(txCtx, user.ID); err != nil {
					return err
				}
			}

			return ErrInvalidCredentials
		}

		response, err = s.AuthenticateUser(txCtx, user, userAgent, ipAddress)

		return err
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) ActivateAccount(ctx context.Context, userID uuid.UUID, req ActivateAccountRequest) error {
	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		user, err := s.users.FindByID(txCtx, userID)
		if err != nil {
			return ErrUnauthorized
		}

		if user.IsVerified {
			return ErrEmailAlreadyVerified
		}

		hash := s.tokenService.Hash(req.Code)
		token, err := s.tokenService.FindByTokenAndTypeAndUser(
			txCtx,
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
			_ = s.tokenService.DeleteByUserAndType(txCtx, user.ID, domainToken.AccountActivation)

			return ErrVerificationCodeExpired
		}

		err = s.users.Verify(txCtx, user.ID)
		if err != nil {
			return err
		}

		err = s.tokenService.Consume(txCtx, token.ID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (*LoginResponse, error) {
	var response *LoginResponse

	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		hashed := s.tokenService.Hash(req.RefreshToken)

		session, err := s.sessions.FindByToken(txCtx, hashed)
		if err != nil {
			return ErrInvalidRefreshToken
		}

		if session.RevokedAt != nil {
			return ErrInvalidRefreshToken
		}

		if session.ExpiresAt.Before(time.Now()) {
			_ = s.sessions.RevokeSession(
				txCtx,
				session.ID,
				backendSession.TokenExpired,
			)
			return ErrSessionExpired
		}

		err = s.sessions.UpdateSessionActivity(txCtx, session.ID)
		if err != nil {
			return ErrInvalidRefreshToken
		}

		user, err := s.users.FindByID(txCtx, session.UserID)
		if err != nil {
			return ErrInvalidRefreshToken
		}

		accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, session.ID)
		if err != nil {
			return err
		}

		refreshToken, err := s.tokenService.CreateRefreshToken()
		if err != nil {
			return err
		}

		err = s.sessions.UpdateRefreshToken(txCtx, session.ID, refreshToken.TokenHash, refreshToken.ExpiresAt)
		if err != nil {
			return err
		}

		response = &LoginResponse{
			User:         userResponse.NewUserResponse(user, s.baseURL),
			AccessToken:  accessToken,
			RefreshToken: refreshToken.Token,
			ExpiresIn:    s.cfg.AccessTokenMinutes,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) Logout(ctx context.Context, userID, sessionID uuid.UUID, ipAddress, userAgent string) error {
	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		err := s.sessions.RevokeSession(
			txCtx,
			sessionID,
			backendSession.RevokedByLogout,
		)
		if err != nil {
			return ErrInvalidRefreshToken
		}

		err = s.sessions.UpdateSessionActivity(txCtx, sessionID)
		if err != nil {
			return ErrInvalidRefreshToken
		}

		deviceInfo := ParseDeviceInfo(userAgent, ipAddress)
		err = s.security.RecordLogout(
			txCtx,
			userID,
			sessionID,
			ipAddress,
			userAgent,
			deviceInfo.DeviceName,
		)
		if err != nil {
			return ErrInvalidRefreshToken
		}

		return nil
	})

	return err
}

func (s *Service) LogoutAllSessions(ctx context.Context, userID uuid.UUID, ipAddress, userAgent string) error {
	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		activeSessions, err := s.sessions.GetSessionsByUserID(txCtx, userID)
		if err != nil {
			return err
		}

		revokeErr := s.sessions.RevokeAllSessions(
			txCtx,
			userID,
			backendSession.RevokedByLogoutAll,
		)
		if revokeErr != nil {
			return ErrInvalidRefreshToken
		}

		for _, session := range activeSessions {
			err = s.sessions.UpdateSessionActivity(txCtx, session.ID)
			if err != nil {
				return err
			}

			err = s.security.RecordLogout(txCtx, session.UserID, session.ID, ipAddress, userAgent, parseDeviceName(userAgent))
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
