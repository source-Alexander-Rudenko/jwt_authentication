package usecase

import (
	"jwt_auth_project/internal/domain"
	"jwt_auth_project/internal/repo"
)

// UserUseCase - интерфейс для бизнес-логики работы с пользователями
type UserUseCase interface {
	GetUserByID(id int64) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	CreateUser(user domain.User) error
}

type userUseCase struct {
	repo repo.UserRepository
}

// NewUserUsecase - конструктор для создания нового userUseCase
func NewUserUsecase(repo repo.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) GetUserByID(id int64) (*domain.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *userUseCase) GetUserByEmail(email string) (*domain.User, error) {
	return u.repo.GetUserByEmail(email)
}

func (u *userUseCase) CreateUser(user domain.User) error {
	return u.repo.CreateUser(user)
}
