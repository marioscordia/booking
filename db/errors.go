package db

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")
	
	ErrDuplicateEmail = errors.New("duplicate email address provided")
)