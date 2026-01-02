package service

import (
	//"time"

	"github.com/RomaNano/go-backend-portfolio-bimbo-version/users-api-http/internal/model"
	"github.com/RomaNano/go-backend-portfolio-bimbo-version/users-api-http/internal/repository"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Create(email string) (*model.User, error) {
	if email == "" {
		return nil, ErrEmailRequired
	}

	user := &model.User{
		Email: email,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}


func (s *UserServiceImpl) List() ([]model.User, error) {
	return []model.User{}, nil
}

func (s *UserServiceImpl) GetByID(id int64) (*model.User, error) {
	return nil, ErrUserNotFound
}
