package strategy

import (
	"errors"
)

var (
	//todo: возможно имеет смысл запилить отдельный тип ошибок, когда ноды не прошли проверку по capabilities
	ErrNotFound      = errors.New("strategy: not found available nodes")
	ErrNotApplicable = errors.New("strategy: not applicable for current node")
)
