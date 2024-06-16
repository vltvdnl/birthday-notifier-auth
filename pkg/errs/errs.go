package errs

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrAppNotFound        = errors.New("app not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrPGNoConnection     = errors.New("pg is not connected")
	ErrPGNoAnswer         = errors.New("pg is not answered")
)
