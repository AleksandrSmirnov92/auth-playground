package domain

// UserRepository — контракт хранилища пользователей.
//
// UseCase зависит от этого интерфейса (а не от конкретной БД),
// поэтому мы можем легко заменить in-memory реализацию на Postgres/Redis.
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Delete(id string) error
}
