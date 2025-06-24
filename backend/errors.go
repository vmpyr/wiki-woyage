package main

import (
	"errors"
)

var (
	ErrUsernameTaken    = errors.New("username already taken")
	ErrConnectionClosed = errors.New("connection closed by client")
)
