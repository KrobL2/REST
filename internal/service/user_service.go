package service

import "go-rest-server/internal/domain"

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListUsers() ([]domain.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) CreateUser(name, email string) (domain.User, error) {
	user := domain.User{
		Name:  name,
		Email: email,
	}

	return s.repo.Create(user)
}
