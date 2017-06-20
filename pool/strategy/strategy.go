package strategy

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("strategy: not found available nodes")
	ErrNotApplicable = errors.New("strategy: not applicable for current node")
)
