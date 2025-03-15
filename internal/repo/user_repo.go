package repo

import (
	"database/sql"
	"errors"
	"jwt/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

var _ domain.UserRepository = (*UserRepo)(nil)

func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	u := new(domain.User)
	for rows.Next() {
		u, err = r.scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *UserRepo) GetUserByID(id int64) (*domain.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	u := new(domain.User)
	for rows.Next() {
		u, err = r.scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (r *UserRepo) CreateUser(user domain.User) error {
	return nil
}

func (r *UserRepo) scanRowIntoUser(rows *sql.Rows) (*domain.User, error) {
	user := new(domain.User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
