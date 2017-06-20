package persistent

import (
	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/qa-dev/jsonwire-grid/pool"
)

type StrategyFactory struct{}

func (f *StrategyFactory) Create(config config.Config, storage pool.StorageInterface) (pool.StrategyInterface, error) {
	return &Strategy{storage}, nil
}
