package entities

import "errors"

var (
	ErrNotValidLogin       = errors.New("user with this login is exists")
	ErrNotValidLoginOrPass = errors.New("not valid login or password")
)
