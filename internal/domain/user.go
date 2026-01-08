package domain

type User struct {
	ID    int
	Name  string
	Email string
}

// Интерфейс репозитория
type UserRepository interface {
	GetAll() ([]User, error)
	Create(user User) (User, error)
}
