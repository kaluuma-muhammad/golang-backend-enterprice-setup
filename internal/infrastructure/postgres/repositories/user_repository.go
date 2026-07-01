package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/go-api/internal/domain/user"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {

	return &UserRepository{
		q: q,
	}
}

func toDomain(u db.User) *user.User {

	return &user.User{
		ID:         u.ID,
		Email:      u.Email,
		Password:   u.Password,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {

	return r.q.CreateUser(
		ctx,
		db.CreateUserParams{
			ID:         u.ID,
			Email:      u.Email,
			Password:   u.Password,
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			IsVerified: u.IsVerified,
			CreatedAt:  u.CreatedAt,
			UpdatedAt:  u.UpdatedAt,
		},
	)
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {

	record, err := r.q.GetUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return toDomain(record), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {

	record, err := r.q.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return toDomain(record), nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {

	return r.q.UpdateUser(
		ctx,
		db.UpdateUserParams{
			ID:         u.ID,
			Email:      u.Email,
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			IsVerified: u.IsVerified,
			UpdatedAt:  u.UpdatedAt,
		},
	)
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {

	return r.q.DeleteUser(ctx, id)
}
