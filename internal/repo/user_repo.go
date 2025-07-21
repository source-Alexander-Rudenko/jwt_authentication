package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"jwt_auth_project/internal/domain"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo создаст объект который, будет удовлетворять
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := new(domain.User)
	err := r.pool.
		QueryRow(ctx,
			`SELECT id, first_name, last_name, email, password, created_at
             FROM users
             WHERE email = $1`, email).
		Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
		)

	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	u := new(domain.User)
	err := r.pool.
		QueryRow(ctx,
			`SELECT id, first_name, last_name, email, password, created_at
             FROM users
             WHERE id = $1`, id).
		Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
		)

	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return u, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, user *domain.User) error {
	err := r.pool.
		QueryRow(ctx,
			`INSERT INTO users (first_name, last_name, email, password)
             VALUES ($1, $2, $3, $4)
             RETURNING id, created_at`,
			user.FirstName,
			user.LastName,
			user.Email,
			user.Password,
		).
		Scan(&user.ID, &user.CreatedAt)

	return err
}
