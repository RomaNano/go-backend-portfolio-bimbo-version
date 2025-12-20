package service

import "users-api-gin/internal/model"

type UserService interface {
	CreateUser(email string) (*model.User, error)
	GetUser(id int64) (*model.User, error)
	ListUsers() ([]model.User, error)
}
