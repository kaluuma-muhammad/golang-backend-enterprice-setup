package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/domain/user"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/go-api/internal/infrastructure/postgres/types"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {

	return &UserRepository{
		pool: pool,
	}
}

func toDomain(u db.User) *user.User {

	return &user.User{
		ID:         u.ID,
		Email:      u.Email,
		Password:   u.Password,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Phone:      types.FromNullablePGText(u.Phone),
		ImageURL:   types.FromNullablePGText(u.ImageUrl),
		IsVerified: u.IsVerified,
		CreatedAt:  types.FromPGTimestamp(u.CreatedAt),
		UpdatedAt:  types.FromPGTimestamp(u.UpdatedAt),
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	q := common.GetQueries(ctx, r.pool)

	return q.CreateUser(
		ctx,
		db.CreateUserParams{
			ID:         u.ID,
			Email:      u.Email,
			Password:   u.Password,
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			IsVerified: u.IsVerified,
			CreatedAt:  types.ToPGTimestamp(u.CreatedAt),
			UpdatedAt:  types.ToPGTimestamp(u.UpdatedAt),
		},
	)
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return toDomain(record), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	q := common.GetQueries(ctx, r.pool)

	record, err := q.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return toDomain(record), nil
}

func (r *UserRepository) UpdateUserAccount(ctx context.Context, u *user.User) error {
	q := common.GetQueries(ctx, r.pool)

	return q.UpdateUserAccount(
		ctx,
		db.UpdateUserAccountParams{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Phone:     types.ToNullablePGText(u.Phone),
		},
	)
}

func (r *UserRepository) UpdateUserAvatar(ctx context.Context, u *user.User) error {
	q := common.GetQueries(ctx, r.pool)

	return q.UpdateUserAvatar(
		ctx,
		db.UpdateUserAvatarParams{
			ID:       u.ID,
			ImageUrl: types.ToNullablePGText(u.ImageURL),
		},
	)
}

func (r *UserRepository) Verify(ctx context.Context, id uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)
	return q.VerifyUser(ctx, id)
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	q := common.GetQueries(ctx, r.pool)
	return q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:       id,
		Password: password,
	})
}

func (r *UserRepository) IncrementFailedLoginAttempts(ctx context.Context, id uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)
	return q.IncrementFailedLoginAttempts(ctx, id)
}

func (r *UserRepository) ResetFailedLoginAttempts(ctx context.Context, id uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)
	return q.ResetFailedLoginAttempts(ctx, id)
}

func (r *UserRepository) LockUserAccount(ctx context.Context, id uuid.UUID, until time.Time) error {
	q := common.GetQueries(ctx, r.pool)
	return q.LockUserAccount(
		ctx,
		db.LockUserAccountParams{
			ID:          id,
			LockedUntil: types.ToPGTimestamp(until),
		},
	)
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)
	return q.UpdateLastLogin(ctx, id)
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := common.GetQueries(ctx, r.pool)
	return q.DeleteUser(ctx, id)
}
