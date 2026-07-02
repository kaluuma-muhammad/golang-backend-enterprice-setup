package repositories

import (
	"context"

	"github.com/go-api/internal/domain/token"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/google/uuid"
)

type TokenRepository struct {
	q *db.Queries
}

func NewTokenRepository(q *db.Queries) *TokenRepository {
	return &TokenRepository{q: q}
}

func toTokenDomain(t db.Token) *token.Token {
	return &token.Token{
		ID:        t.ID,
		UserID:    t.UserID,
		Type:      token.Type(t.Type),
		Token:     t.Token,
		ExpiresAt: t.ExpiresAt,
		UsedAt:    &t.UsedAt,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func (r *TokenRepository) Create(ctx context.Context, token *token.Token) error {
	return r.q.CreateToken(ctx, db.CreateTokenParams{
		ID:        token.ID,
		UserID:    token.UserID,
		Type:      db.TokenType(token.Type),
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		UsedAt:    *token.UsedAt,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	})
}

func (r *TokenRepository) FindByID(ctx context.Context, id uuid.UUID) (*token.Token, error) {
	t, err := r.q.FindByTokenByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toTokenDomain(t), nil
}

func (r *TokenRepository) FindByToken(ctx context.Context, token string) (*token.Token, error) {
	t, err := r.q.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return toTokenDomain(t), nil
}

func (r *TokenRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	return r.q.MarkTokenAsUsed(ctx, id)
}

func (r *TokenRepository) DeleteExpired(ctx context.Context) error {
	return r.q.DeleteExpiredTokens(ctx)
}
