package kubernetes

import (
	"errors"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/qa-dev/jsonwire-grid/pool/strategy"
	"github.com/satori/go.uuid"
	"time"
	"fmt"
)

type Strategy struct {
	storage        pool.StorageInterface
	provider       kubernetesProviderInterface
	config         strategyConfig
	capsComparator capabilities.ComparatorInterface
}

const cleanupFailedPodsTimeout = time.Minute * 2

func (s *Strategy) Reserve(desiredCaps capabilities.Capabilities) (pool.Node, error) {
	nodeConfig := s.findApplicableConfig(s.config.NodeList, desiredCaps)
	if nodeConfig == nil {
		return pool.Node{}, strategy.ErrNotFound
	}
	podName := "wd-node-" + uuid.NewV4().String()
	ts := time.Now().Unix()
	node := pool.NewNode(podName, pool.NodeTypeKubernetes, "temp-value-replace-me", pool.NodeStatusReserved, "", ts, ts, []capabilities.Capabilities{})
	err := s.storage.Add(*node, s.config.Limit)
	if err != nil {
		return pool.Node{}, errors.New("add node to storage, " + err.Error())
	}
	nodeAddress, err := s.provider.Create(podName, nodeConfig.Params)
	if err != nil {
		go func(podName string) {
			time.Sleep(cleanupFailedPodsTimeout)
			_ = s.provider.Destroy(podName) // на случай если что то криво создалось
		}(podName)
		return pool.Node{}, errors.New("create node by provider, " + err.Error())
	}

	err = s.storage.UpdateAddress(*node, nodeAddress)
	if err != nil {
		return pool.Node{}, errors.New("update node address in storage, " + err.Error())
	}
	node.Address = nodeAddress
	return *node, nil

}

func (s *Strategy) CleanUp(node pool.Node) error {
	if node.Type != pool.NodeTypeKubernetes {
		return strategy.ErrNotApplicable
	}
	if node.Key == "" {
		return fmt.Errorf("empty node key")
	}
	err := s.provider.Destroy(node.Key)
	if err != nil {
		return errors.New("destroy node by provider, " + err.Error())
	}
	err = s.storage.Remove(node)
	if err != nil {
		return errors.New("remove node from storage, " + err.Error())
	}
	return nil
}

func (s *Strategy) FixNodeStatus(node pool.Node) error {
	if node.Type != pool.NodeTypeKubernetes {
		return strategy.ErrNotApplicable
	}
	if node.Key == "" {
		return fmt.Errorf("empty node key")
	}
	err := s.provider.Destroy(node.Key)
	if err != nil {
		return errors.New("destroy node by provider, " + err.Error())
	}
	err = s.storage.Remove(node)
	if err != nil {
		return errors.New("remove node from storage, " + err.Error())
	}
	return nil
}

// findApplicableConfig смотрим в конфиг стратегии, есть ли там подходящие ноды
func (s *Strategy) findApplicableConfig(nodeList []nodeConfig, desiredCaps capabilities.Capabilities) *nodeConfig {
	for _, nodeConfig := range nodeList {
		for _, nodeCapabilities := range nodeConfig.CapabilitiesList {
			if s.capsComparator.Compare(desiredCaps, nodeCapabilities) {
				return &nodeConfig
			}
		}
	}
	return nil
}
