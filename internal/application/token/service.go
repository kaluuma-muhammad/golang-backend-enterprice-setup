package token

import (
	"context"
	"time"

	domainToken "github.com/go-api/internal/domain/token"

	"github.com/google/uuid"
)

type Service struct {
	repository domainToken.Repository
	generator  *Generator
}

const RefreshTokenTTL = 7 * 24 * time.Hour

func NewService(repository domainToken.Repository, generator *Generator) *Service {
	return &Service{
		repository: repository,
		generator:  generator,
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
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
	}, nil
}

func (s *Service) Hash(token string) string {

	return s.generator.Hash(token)
}
