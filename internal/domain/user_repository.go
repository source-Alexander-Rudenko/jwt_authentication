package domain

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	CreateUser(user User) error
}
