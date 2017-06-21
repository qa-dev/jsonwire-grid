package pool

import (
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (s *StorageMock) Add(node Node, limit int) error {
	args := s.Called(node, limit)
	return args.Error(0)
}

func (s *StorageMock) ReserveAvailable(nodeList []Node) (Node, error) {
	args := s.Called(nodeList)
	return args.Get(0).(Node), args.Error(1)
}

func (s *StorageMock) SetBusy(node Node, sessionId string) error {
	args := s.Called(node, sessionId)
	return args.Error(0)
}

func (s *StorageMock) SetAvailable(node Node) error {
	args := s.Called(node)
	return args.Error(0)
}

func (s *StorageMock) GetCountWithStatus(nodeStatus *NodeStatus) (int, error) {
	args := s.Called(nodeStatus)
	return args.Int(0), args.Error(1)
}

func (s *StorageMock) GetBySession(sessionId string) (Node, error) {
	args := s.Called(sessionId)
	return args.Get(0).(Node), args.Error(1)
}

func (s *StorageMock) GetByAddress(address string) (Node, error) {
	args := s.Called(address)
	return args.Get(0).(Node), args.Error(1)
}

func (s *StorageMock) GetAll() ([]Node, error) {
	args := s.Called()
	return args.Get(0).([]Node), args.Error(1)
}

func (s *StorageMock) Remove(node Node) error {
	args := s.Called(node)
	return args.Error(0)
}

type StrategyListMock struct {
	mock.Mock
}

func (s *StrategyListMock) Reserve(caps capabilities.Capabilities) (Node, error) {
	args := s.Called(caps)
	return args.Get(0).(Node), args.Error(1)
}

func (s *StrategyListMock) CleanUp(node Node) error {
	args := s.Called(node)
	return args.Error(0)
}

func (s *StrategyListMock) FixNodeStatus(node Node) error {
	args := s.Called(node)
	return args.Error(0)
}
