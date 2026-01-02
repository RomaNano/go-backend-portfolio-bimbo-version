package service

import (
	"github.com/RomaNano/go-backend-portfolio-bimbo-version/users-api-http/internal/model"
)

type UserService interface{
	Create(email string) (*model.User, error)
	List()([]model.User, error)
	GetByID(id int64) (*model.User, error)
}