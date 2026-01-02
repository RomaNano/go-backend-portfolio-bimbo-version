package repository

import (
	"github.com/romanrey/go-backend-portfolio-bimbo-version/users-api-http/internal/model"
)


type UserRepository interface{
	Create(user *model.User) error
	GetByID(id int64) (*model.User, error)
	List() ([]model.User, error)
}
