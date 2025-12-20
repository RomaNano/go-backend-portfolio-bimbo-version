package service

import (
	"users-api-gin/internal/model"
	"users-api-gin/internal/repository"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}


func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

// CreateUser
func (s *UserServiceImpl) CreateUser(email string) (*model.User, error) {
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

// GetUser
func (s *UserServiceImpl) GetUser(id int64) (*model.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// ListUsers
func (s *UserServiceImpl) ListUsers() ([]model.User, error) {
	return s.repo.List()
}