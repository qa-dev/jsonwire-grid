package persistent

import (
	"errors"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/strategy"
)

type Strategy struct {
	storage pool.StorageInterface
}

func (s *Strategy) Reserve(capabilities pool.Capabilities) (pool.Node, error) {
	// todo: выпилить логику сравнения capabilities из storage в какой-нибудь CapabilitiesComparator
	node, err := s.storage.ReserveAvailable(capabilities)
	if err != nil {
		return pool.Node{}, errors.New("Get persistent node, " + err.Error())
	}
	return node, err
}

func (s *Strategy) CleanUp(node pool.Node) error {
	if node.Type != pool.NodeTypePersistent {
		return strategy.ErrNotApplicable
	}
	err := s.storage.SetAvailable(node)
	if err != nil {
		return errors.New("CleanUp persistent node, " + err.Error())
	}
	return nil
}
