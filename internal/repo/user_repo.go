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
	CreateUser(ctx context.Context, user domain.User) error
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := new(domain.User)
	err := r.pool.
		QueryRow(ctx,
			`SELECT id, username, email, password_hash, created_at
             FROM "USER"
             WHERE email = $1`, email).
		Scan(
			&u.ID,
			&u.Username,
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
			`SELECT id, username, email, password_hash, created_at
             FROM "USER"
             WHERE id = $1`, id).
		Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
		)

	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return u, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, user domain.User) error {
	_, err := r.pool.
		Exec(ctx,
			`INSERT INTO "USER" (id, username, email, password_hash, created_at)
             VALUES ($1, $2, $3, $4, $5)
             `,
			user.ID,
			user.Username,
			user.Email,
			user.Password,
			user.CreatedAt,
		)

	return err
}
