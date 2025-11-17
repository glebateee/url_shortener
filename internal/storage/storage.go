package storage

import "errors"

var (
	ErrURLExists   = errors.New("alias already exists")
	ErrURLNotFound = errors.New("alias not found")
)
