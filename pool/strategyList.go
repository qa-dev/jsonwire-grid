package pool

import (
	"github.com/qa-dev/jsonwire-grid/pool/strategy"
	"errors"
)

type StrategyListInterface interface {
	Reserve(caps Capabilities) (node Node, err error)
	CleanUp(node Node) error
}

type StrategyList struct {
	list []StrategyInterface
}

func NewStrategyList(list []StrategyInterface) *StrategyList {
	return &StrategyList{list}
}

func (s *StrategyList) Reserve(caps Capabilities) (node Node, err error) {
	err = errors.New("Empty strategy list")
	for _, currStrategy := range s.list {
		node, err = currStrategy.Reserve(caps)
		if err == strategy.ErrNotFound {
			continue
		}
		break
	}
	return
}

func (s *StrategyList) CleanUp(node Node) error {
	err := errors.New("Empty strategy list")
	for _, currStrategy := range s.list {
		err = currStrategy.CleanUp(node)
		if err == strategy.ErrNotApplicable {
			continue
		}
		break
	}
	return err
}