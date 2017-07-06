package persistent

import (
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
)

type StrategyFactory struct{}

func (f *StrategyFactory) Create(
	storage pool.StorageInterface,
	capsComparator capabilities.ComparatorInterface,
	clientFactory jsonwire.ClientFactoryInterface,
) (pool.StrategyInterface, error) {
	return &Strategy{
		storage,
		capsComparator,
		clientFactory,
		new(nodeHelperFactory),
	}, nil
}
