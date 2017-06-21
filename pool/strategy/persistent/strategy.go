package persistent

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/qa-dev/jsonwire-grid/pool/strategy"
)

type Strategy struct {
	storage        pool.StorageInterface
	capsComparator capabilities.ComparatorInterface
	clientFactory  jsonwire.ClientFactoryInterface
}

func (s *Strategy) Reserve(desiredCaps capabilities.Capabilities) (pool.Node, error) {
	//todo: -- begin подумать над тем чтобы выполнять эти задачи в фоне
	nodeList, err := s.storage.GetAll()
	if err != nil {
		return pool.Node{}, errors.New("Get all desiredCpos list, " + err.Error())
	}
	s.registerCapabilities(nodeList)
	// todo: -- end

	applicableNodeList := s.findApplicableNodes(nodeList, desiredCaps)

	// цикл для того чтобы не уйти в рекурсию, в случае когда все ноды не работают, но регистрируются быстрее чем выпиливаются
	for i := range applicableNodeList {
		node, err := s.storage.ReserveAvailable(applicableNodeList[i:])
		if err != nil {
			log.Errorf("reserve node in storage, %s", err)
			break
		}
		client := s.clientFactory.Create(node.Address)
		message, err := client.Status()
		if err != nil {
			log.Infof("Remove unavailable node [%s], %s", node, err)
			err = s.storage.Remove(node)
			if err != nil {
				log.Errorf("Remove unavailable node [%s], %s", node, err)
			}
			continue
		}
		//todo: заменить магические числа на константы статусов
		if message.Status == 0 { // status == ok
			return node, nil
		}
	}

	return pool.Node{}, strategy.ErrNotFound
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

func (s *Strategy) FixNodeStatus(node pool.Node) error {
	if node.Type != pool.NodeTypePersistent {
		return strategy.ErrNotApplicable
	}
	err := s.storage.SetAvailable(node)
	if err != nil {
		return errors.New("fix node status to available, " + err.Error())
	}
	return nil
}

func (s *Strategy) findApplicableNodes(nodeList []pool.Node, desiredCaps capabilities.Capabilities) []pool.Node {
	var applicableNodeList []pool.Node
	for _, node := range nodeList {
		for _, availableCaps := range node.CapabilitiesList {
			if s.capsComparator.Compare(desiredCaps, availableCaps) {
				applicableNodeList = append(applicableNodeList, node)
			}
		}
	}
	return applicableNodeList
}

func (s *Strategy) registerCapabilities(nodeList []pool.Node) {
	for _, node := range nodeList {
		for _, availableCaps := range node.CapabilitiesList {
			s.capsComparator.Register(availableCaps)

		}
	}
}
