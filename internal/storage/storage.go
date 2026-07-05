// Package storage defines signal errors shared by all storages.
package storage

import (
	"errors"
)

var (
	// ErrNotFound is returned by storage when code/URL is not found.
	ErrNotFound = errors.New("not found")
	// ErrCodeTaken is returned by storage when the code is already taken by another URL.
	ErrCodeTaken = errors.New("code already taken")
	// ErrURLExists is returned by storage when the URL already has a code.
	ErrURLExists = errors.New("url already exists")
)
