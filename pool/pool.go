package pool

import (
	"time"

	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
)

const (
	// todo: сейчас не используются, совсем удалить или сделать дефолтные параметры для конфига.
	defaultBusyNodeDuration     = time.Minute * 30
	defaultReservedNodeDuration = time.Minute * 5
)

type StorageInterface interface {
	Add(node Node, limit int) error
	ReserveAvailable([]Node) (Node, error)
	SetBusy(Node, string) error
	SetAvailable(Node) error
	GetCountWithStatus(*NodeStatus) (int, error)
	GetBySession(string) (Node, error)
	GetByAddress(string) (Node, error)
	GetAll() ([]Node, error)
	Remove(Node) error
}

type StrategyInterface interface {
	Reserve(capabilities.Capabilities) (Node, error)
	CleanUp(Node) error
	FixNodeStatus(Node) error
}

type Pool struct {
	storage              StorageInterface
	busyNodeDuration     time.Duration
	reservedNodeDuration time.Duration
	strategyList         StrategyListInterface
}

func NewPool(storage StorageInterface, strategyList StrategyListInterface) *Pool {
	return &Pool{
		storage:              storage,
		busyNodeDuration:     defaultBusyNodeDuration,
		reservedNodeDuration: defaultReservedNodeDuration,
		strategyList:         strategyList,
	}
}

func (p *Pool) SetBusyNodeDuration(duration time.Duration) {
	p.busyNodeDuration = duration
}

func (p *Pool) SetReservedNodeDuration(duration time.Duration) {
	p.reservedNodeDuration = duration
}

// TODO: research close transaction and defer close mysql result body.
func (p *Pool) ReserveAvailableNode(caps capabilities.Capabilities) (*Node, error) {
	node, err := p.strategyList.Reserve(caps)
	if err != nil {
		err = errors.New("Can't reserve available node, " + err.Error())
		log.Error(err)
		return nil, err
	}
	return &node, err
}

func (p *Pool) Add(key string, t NodeType, address string, capabilitiesList []capabilities.Capabilities) error {
	if len(capabilitiesList) == 0 {
		return errors.New("[Pool/Add] Capabilities must contains more one element")
	}
	ts := time.Now().Unix()
	return p.storage.Add(*NewNode(key, t, address, NodeStatusAvailable, "", ts, ts, capabilitiesList), 0)
}

func (p *Pool) RegisterSession(node *Node, sessionID string) error {
	return p.storage.SetBusy(*node, sessionID)
}

func (p *Pool) GetAll() ([]Node, error) {
	return p.storage.GetAll()
}

func (p *Pool) GetNodeBySessionID(sessionID string) (*Node, error) {
	node, err := p.storage.GetBySession(sessionID)
	if err != nil {
		err = fmt.Errorf("Can't find node by session[%s], %s", sessionID, err.Error())
		log.Error(err)
		return nil, err
	}
	return &node, nil
}

func (p *Pool) GetNodeByAddress(address string) (*Node, error) {
	node, err := p.storage.GetByAddress(address)
	if err != nil {
		err = fmt.Errorf("Can't find node by address[%s], %s", address, err.Error())
		log.Error(err)
		return nil, err
	}
	return &node, nil
}

func (p *Pool) CleanUpNode(node *Node) error {
	err := p.strategyList.CleanUp(*node)
	if err != nil {
		err = errors.New("Can't clean up node: " + err.Error())
		log.Error(err)
	}
	return err
}

// удаляет ноду из пула
func (p *Pool) Remove(node *Node) error {
	err := p.storage.Remove(*node)
	if err != nil {
		err = errors.New("Can't remove node from pool, " + err.Error())
		log.Error(err)
		return err
	}
	return nil
}

func (p *Pool) CountNodes(status *NodeStatus) (int, error) {
	count, err := p.storage.GetCountWithStatus(status)
	if err != nil {
		err = errors.New("Can't get count nodes, " + err.Error())
		log.Error(err)
		return 0, err
	}
	return count, nil
}

func (p *Pool) FixNodeStatuses() {
	nodeList, err := p.GetAll()
	if err != nil {
		err = errors.New("Can't check node statuses, " + err.Error())
		log.Error(err)
		return
	}
	for _, node := range nodeList {
		isFixed, err := p.fixNodeStatus(&node)
		if err != nil {
			log.Error(err)
			continue
		}
		if isFixed {
			log.Infof("Node [%s] status fixed", node.Key)
		}
	}
}

func (p *Pool) fixNodeStatus(node *Node) (bool, error) {
	nodeStatusDuration := time.Since(time.Unix(node.Updated, 0))
	switch node.Status {
	case NodeStatusReserved:
		if nodeStatusDuration < p.reservedNodeDuration {
			return false, nil
		}
	case NodeStatusBusy:
		if nodeStatusDuration < p.busyNodeDuration {
			return false, nil
		}
	default:
		return false, nil
	}
	err := p.strategyList.FixNodeStatus(*node)
	if err != nil {
		return false, fmt.Errorf("Can't fix node [%s] status, %s", node.Key, err.Error())
	}
	return true, nil
}
