package storage

import (
	"errors"
)

var (
	ErrNotFound = errors.New("storage: not found available nodes")
)
