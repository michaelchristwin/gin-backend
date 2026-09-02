package domain

import (
	"errors"
)

var (
	// Auth errors
	ErrEmailTaken         = errors.New("Email taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordTooWeak    = errors.New("password does not meet requirements")
	ErrPasswordPwned      = errors.New("this password has appeared in a known data breach — please choose a different one")

	// User errors
	ErrUserNotFound = errors.New("user not found")
)
