package user

import "errors"

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
)
