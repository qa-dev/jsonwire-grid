package local

import (
	"errors"
	"sync"
	"time"

	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/storage"
)

type Storage struct {
	mu sync.RWMutex
	db map[string]*pool.Node
}

func (s *Storage) Add(node pool.Node, limit int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit > 0 {
		i := 0
		for _, currNode := range s.db {
			if currNode.Status == node.Status {
				i++
			}
		}
		if i >= limit {
			return errors.New("limit reached")
		}
	}

	s.db[node.Key] = &node
	return nil
}

func (s *Storage) ReserveAvailable(nodeList []pool.Node) (pool.Node, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, node := range nodeList {
		dbNode, ok := s.db[node.Key]
		if ok && dbNode.Status == pool.NodeStatusAvailable {
			dbNode.Status = pool.NodeStatusReserved
			dbNode.Updated = time.Now().Unix()
			return *dbNode, nil
		}
	}
	return pool.Node{}, storage.ErrNotFound
}

func (s *Storage) SetBusy(node pool.Node, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	storedNode, ok := s.db[node.Key]
	if !ok {
		return storage.ErrNotFound
	}
	storedNode.Status = pool.NodeStatusBusy
	storedNode.SessionID = sessionID
	storedNode.Updated = time.Now().Unix()
	return nil
}

func (s *Storage) SetAvailable(node pool.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	storedNode, ok := s.db[node.Key]
	if !ok {
		return storage.ErrNotFound
	}
	storedNode.Status = pool.NodeStatusAvailable
	storedNode.Updated = time.Now().Unix()
	return nil
}

func (s *Storage) GetCountWithStatus(status *pool.NodeStatus) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if status == nil {
		return len(s.db), nil
	}
	count := 0
	for _, node := range s.db {
		if node.Status == *status {
			count++
		}
	}
	return count, nil
}

func (s *Storage) GetBySession(sessionID string) (pool.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, node := range s.db {
		if node.SessionID == sessionID {
			return *node, nil
		}
	}
	return pool.Node{}, storage.ErrNotFound
}

func (s *Storage) GetByAddress(address string) (pool.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, node := range s.db {
		if node.Address == address {
			return *node, nil
		}
	}
	return pool.Node{}, storage.ErrNotFound
}

func (s *Storage) GetAll() ([]pool.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	nodeList := make([]pool.Node, 0, len(s.db))
	for _, value := range s.db {
		nodeList = append(nodeList, *value)
	}

	return nodeList, nil
}

func (s *Storage) Remove(node pool.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.db[node.Key]
	if !ok {
		return storage.ErrNotFound
	}
	delete(s.db, node.Key)
	return nil
}

func (s *Storage) UpdateAddress(node pool.Node, newAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	storedNode, ok := s.db[node.Key]
	if !ok {
		return storage.ErrNotFound
	}
	storedNode.Address = newAddress
	return nil
}
