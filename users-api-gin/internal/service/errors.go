package service

import "errors"

var (
	ErrEmailRequired     = errors.New("email is required")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
