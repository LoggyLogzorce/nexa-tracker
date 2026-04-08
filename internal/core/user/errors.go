package user

import "errors"

var (
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserOwnsProjects = errors.New("user owns projects")
	ErrUserNotFound     = errors.New("user not found")
)
