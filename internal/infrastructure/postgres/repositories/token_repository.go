package repositories

import (
	"context"

	domainToken "github.com/go-api/internal/domain/token"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/go-api/internal/infrastructure/postgres/types"
	"github.com/google/uuid"
)

type TokenRepository struct {
	q *db.Queries
}

func NewTokenRepository(q *db.Queries) *TokenRepository {
	return &TokenRepository{q: q}
}

func toTokenDomain(t db.Token) *domainToken.Token {
	return &domainToken.Token{
		ID:        t.ID,
		UserID:    t.UserID,
		Type:      domainToken.Type(t.Type),
		Token:     t.Token,
		ExpiresAt: types.FromPGTimestamp(t.ExpiresAt),
		UsedAt:    types.FromPGTimestampPtr(t.UsedAt),
		CreatedAt: types.FromPGTimestamp(t.CreatedAt),
		UpdatedAt: types.FromPGTimestamp(t.UpdatedAt),
	}
}

func (r *TokenRepository) Create(ctx context.Context, token *domainToken.Token) error {
	return r.q.CreateToken(ctx, db.CreateTokenParams{
		ID:        token.ID,
		UserID:    token.UserID,
		Type:      db.TokenType(token.Type),
		Token:     token.Token,
		ExpiresAt: types.ToPGTimestamp(token.ExpiresAt),
		CreatedAt: types.ToPGTimestamp(token.CreatedAt),
		UpdatedAt: types.ToPGTimestamp(token.UpdatedAt),
	})
}

func (r *TokenRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainToken.Token, error) {
	t, err := r.q.FindByTokenByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toTokenDomain(t), nil
}

func (r *TokenRepository) FindByToken(ctx context.Context, token string) (*domainToken.Token, error) {
	t, err := r.q.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return toTokenDomain(t), nil
}

func (r *TokenRepository) FindByTokenAndType(ctx context.Context, token string, tokenType domainToken.Type) (*domainToken.Token, error) {
	record, err := r.q.FindByTokenAndType(ctx, db.FindByTokenAndTypeParams{
		Token: token,
		Type:  db.TokenType(tokenType),
	})

	if err != nil {
		return nil, err
	}

	return toTokenDomain(record), nil
}

func (r *TokenRepository) FindByTokenAndTypeAndUser(ctx context.Context, token string, tokenType domainToken.Type, userID uuid.UUID) (*domainToken.Token, error) {
	record, err := r.q.FindByTokenAndTypeAndUser(ctx, db.FindByTokenAndTypeAndUserParams{
		Token:  token,
		Type:   db.TokenType(tokenType),
		UserID: userID,
	})

	if err != nil {
		return nil, err
	}

	return toTokenDomain(record), nil
}

func (r *TokenRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	return r.q.MarkTokenAsUsed(ctx, id)
}

func (r *TokenRepository) FindByUserAndType(ctx context.Context, userID uuid.UUID, tokenType domainToken.Type) (*domainToken.Token, error) {
	record, err := r.q.GetTokenByUserAndType(ctx, db.GetTokenByUserAndTypeParams{
		UserID: userID,
		Type:   db.TokenType(tokenType),
	})

	if err != nil {
		return nil, err
	}

	return toTokenDomain(record), nil
}

func (r *TokenRepository) DeleteByUserAndType(ctx context.Context, userID uuid.UUID, tokenType domainToken.Type) error {
	return r.q.DeleteTokenByUserAndType(ctx, db.DeleteTokenByUserAndTypeParams{
		UserID: userID,
		Type:   db.TokenType(tokenType),
	})
}

func (r *TokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteToken(ctx, id)
}

func (r *TokenRepository) DeleteExpired(ctx context.Context) error {
	return r.q.DeleteExpiredTokens(ctx)
}
