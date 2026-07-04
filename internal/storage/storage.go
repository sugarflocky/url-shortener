package storage

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrCodeTaken = errors.New("code already taken")
	ErrURLExists = errors.New("url already exists")
)
