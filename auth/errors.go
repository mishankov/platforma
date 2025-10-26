package auth

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrWrongUserOrPassword = errors.New("wrong user or password")

	ErrInvalidUsername = errors.New("invalid username")
	ErrShortUsername   = errors.New("short username")
	ErrLongUsername    = errors.New("long username")

	ErrInvalidPassword          = errors.New("invalid password")
	ErrShortPassword            = errors.New("short password")
	ErrLongPassword             = errors.New("long password")
	ErrCurrentPasswordIncorrect = errors.New("current password is incorrect")
)
