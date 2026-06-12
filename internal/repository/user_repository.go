package repository

import (
	"context"
	"time"

	sqlcdb "github.com/yourusername/user-api/db/sqlc"
)

// UserRepository defines all persistence operations for the user domain.
type UserRepository interface {
	Create(ctx context.Context, name string, dob time.Time) (int64, error)
	GetByID(ctx context.Context, id uint32) (sqlcdb.User, error)
	Update(ctx context.Context, id uint32, name string, dob time.Time) error
	Delete(ctx context.Context, id uint32) error
	List(ctx context.Context, limit, offset int32) ([]sqlcdb.User, error)
	Count(ctx context.Context) (int64, error)
}

// userRepo is the concrete implementation backed by SQLC-generated queries.
type userRepo struct {
	q sqlcdb.Querier
}

// NewUserRepository returns a UserRepository wired to the given SQLC Querier.
func NewUserRepository(q sqlcdb.Querier) UserRepository {
	return &userRepo{q: q}
}

func (r *userRepo) Create(ctx context.Context, name string, dob time.Time) (int64, error) {
	return r.q.CreateUser(ctx, sqlcdb.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

func (r *userRepo) GetByID(ctx context.Context, id uint32) (sqlcdb.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *userRepo) Update(ctx context.Context, id uint32, name string, dob time.Time) error {
	return r.q.UpdateUser(ctx, sqlcdb.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  dob,
	})
}

func (r *userRepo) Delete(ctx context.Context, id uint32) error {
	return r.q.DeleteUser(ctx, id)
}

func (r *userRepo) List(ctx context.Context, limit, offset int32) ([]sqlcdb.User, error) {
	return r.q.ListUsers(ctx, sqlcdb.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *userRepo) Count(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}
