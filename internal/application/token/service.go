package token

import (
	"context"
	"time"

	domainToken "github.com/go-api/internal/domain/token"
	"github.com/go-api/internal/shared/config"

	"github.com/google/uuid"
)

type Service struct {
	repository domainToken.Repository
	generator  *Generator
	cfg        *config.JWTConfig
}

func NewService(repository domainToken.Repository, generator *Generator, cfg *config.JWTConfig) *Service {
	return &Service{
		repository: repository,
		generator:  generator,
		cfg:        cfg,
	}
}

func (s *Service) Create(ctx context.Context, req CreateTokenRequest) (string, error) {
	token, err := s.generator.Generate()
	if err != nil {
		return "", err
	}

	hash := s.generator.Hash(token)

	err = s.repository.Create(
		ctx,
		&domainToken.Token{
			ID:        uuid.New(),
			UserID:    req.UserID,
			Type:      req.Type,
			Token:     hash,
			ExpiresAt: time.Now().Add(req.TTL),
			CreatedAt: time.Now(),
		},
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Verify(ctx context.Context, req VerifyTokenRequest) (*domainToken.Token, error) {
	hash := s.generator.Hash(req.Token)

	token, err := s.repository.FindByToken(ctx, hash)
	if err != nil {
		return nil, err
	}

	if token.Type != req.Type {
		return nil, ErrInvalidToken
	}

	if token.UsedAt != nil {
		return nil, ErrTokenAlreadyUsed
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return token, nil
}

func (s *Service) Consume(ctx context.Context, id uuid.UUID) error {
	return s.repository.MarkUsed(ctx, id)
}

func (s *Service) Cleanup(ctx context.Context) error {
	return s.repository.DeleteExpired(ctx)
}

func (s *Service) CreateRefreshToken() (*RefreshTokenResult, error) {
	token, err := s.generator.Generate()
	if err != nil {
		return nil, err
	}

	hash := s.generator.Hash(token)

	return &RefreshTokenResult{
		Token:     token,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.RefreshTokenDays) * 24 * time.Hour),
	}, nil
}

func (s *Service) CreateResetToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := s.generator.Generate()
	if err != nil {
		return "", err
	}

	hash := s.generator.Hash(token)

	err = s.repository.Create(
		ctx,
		&domainToken.Token{
			ID:        uuid.New(),
			UserID:    userID,
			Type:      domainToken.PasswordResetGrant,
			Token:     hash,
			ExpiresAt: time.Now().Add(15 * time.Minute),
			CreatedAt: time.Now(),
		},
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Hash(token string) string {

	return s.generator.Hash(token)
}

func (s *Service) CreateCodeToken(ctx context.Context, userID uuid.UUID, tokenType domainToken.Type) (string, error) {
	code, err := s.generator.GenerateCode(6)
	if err != nil {
		return "", err
	}

	hash := s.generator.Hash(code)

	_ = s.repository.DeleteByUserAndType(ctx, userID, tokenType)

	err = s.repository.Create(ctx, &domainToken.Token{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      tokenType,
		Token:     hash,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
	},
	)

	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *Service) DeleteByUserAndType(ctx context.Context, userID uuid.UUID, tokenType domainToken.Type) error {
	return s.repository.DeleteByUserAndType(ctx, userID, tokenType)
}

func (s *Service) FindByTokenAndType(ctx context.Context, token string, tokenType domainToken.Type) (*domainToken.Token, error) {
	return s.repository.FindByTokenAndType(ctx, token, tokenType)
}

func (s *Service) FindByTokenAndTypeAndUser(ctx context.Context, token string, tokenType domainToken.Type, userID uuid.UUID) (*domainToken.Token, error) {
	return s.repository.FindByTokenAndTypeAndUser(ctx, token, tokenType, userID)
}

func (s *Service) FindByUserAndType(ctx context.Context, userID uuid.UUID, tokenType domainToken.Type) (*domainToken.Token, error) {
	return s.repository.FindByUserAndType(ctx, userID, tokenType)
}
