package storage

import (
	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/qa-dev/jsonwire-grid/pool"
)

type StorageFactoryInterface interface {
	Create(config.Config) (pool.StorageInterface, error)
}
