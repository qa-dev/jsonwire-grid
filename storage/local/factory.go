package local

import (
	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/qa-dev/jsonwire-grid/pool"
)

type Factory struct {
}

func (f *Factory) Create(cfg config.Config) (pool.StorageInterface, error) {
	return &Storage{db: make(map[string]*pool.Node)}, nil
}
