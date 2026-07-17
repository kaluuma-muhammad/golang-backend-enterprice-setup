package repositories

import (
	"context"

	"github.com/go-api/internal/application/common"
	domainToken "github.com/go-api/internal/domain/token"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/go-api/internal/infrastructure/postgres/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{
		pool: pool,
	}
}

func (r *TokenRepository) Create(ctx context.Context, token *domainToken.Token) error {
	q := common.GetQueries(ctx, r.pool)

	return q.CreateToken(ctx, db.CreateTokenParams{
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
	q := common.GetQueries(ctx, r.pool)

	t, err := q.FindByTokenByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toTokenDomain(t), nil
}

func (r *TokenRepository) FindByToken(ctx context.Context, token string) (*domainToken.Token, error) {
	q := common.GetQueries(ctx, r.pool)

	t, err := q.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return toTokenDomain(t), nil
}

func (r *TokenRepository) FindByTokenAndType(ctx context.Context, token string, tokenType domainToken.Type) (*domainToken.Token, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.FindByTokenAndType(ctx, db.FindByTokenAndTypeParams{
		Token: token,
		Type:  db.TokenType(tokenType),
	})

	if err != nil {
		return nil, err
	}

	return toTokenDomain(record), nil
}

func (r *TokenRepository) FindByTokenAndTypeAndUser(ctx context.Context, token string, tokenType domainToken.Type, userID uuid.UUID) (*domainToken.Token, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.FindByTokenAndTypeAndUser(ctx, db.FindByTokenAndTypeAndUserParams{
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
	q := common.GetQueries(ctx, r.pool)
	return q.MarkTokenAsUsed(ctx, id)
}

func (r *TokenRepository) FindByUserAndType(ctx context.Context, userID uuid.UUID, tokenType domainToken.Type) (*domainToken.Token, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetTokenByUserAndType(ctx, db.GetTokenByUserAndTypeParams{
		UserID: userID,
		Type:   db.TokenType(tokenType),
	})

	if err != nil {
		return nil, err
	}

	return toTokenDomain(record), nil
}

func (r *TokenRepository) DeleteByUserAndType(ctx context.Context, userID uuid.UUID, tokenType domainToken.Type) error {
	q := common.GetQueries(ctx, r.pool)
	return q.DeleteTokenByUserAndType(ctx, db.DeleteTokenByUserAndTypeParams{
		UserID: userID,
		Type:   db.TokenType(tokenType),
	})
}

func (r *TokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)
	return q.DeleteToken(ctx, id)
}

func (r *TokenRepository) DeleteExpired(ctx context.Context) error {
	q := common.GetQueries(ctx, r.pool)
	return q.DeleteExpiredTokens(ctx)
}
