package repository

import "users-api-gin/internal/model"

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id int64) (*model.User, error)
	List() ([]model.User, error)
}