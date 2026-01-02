package service

import (
	"errors"
)

var (
	ErrEmailRequired  = errors.New("email is required")
	ErrUserNotFound   = errors.New("user not found")
	ErrUserExists     = errors.New("user already exists")
)