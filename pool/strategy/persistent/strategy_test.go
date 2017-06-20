package persistent

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/qa-dev/jsonwire-grid/pool"
	"time"
	"errors"
)

type StorageMock struct {
	mock.Mock
}

func (s *StorageMock) Add(pool.Node) error {
	panic("implement me")
}

func (s *StorageMock) ReserveAvailable(caps pool.Capabilities) (pool.Node, error) {
	args := s.Called(caps)
	return args.Get(0).(pool.Node), args.Error(1)
}

func (s *StorageMock) SetBusy(pool.Node, string) error {
	panic("implement me")
}

func (s *StorageMock) SetAvailable(node pool.Node) error {
	args := s.Called(node)
	return args.Error(0)
}

func (s *StorageMock) GetCountWithStatus(*pool.NodeStatus) (int, error) {
	panic("implement me")
}

func (s *StorageMock) GetBySession(string) (pool.Node, error) {
	panic("implement me")
}

func (s *StorageMock) GetByAddress(string) (pool.Node, error) {
	panic("implement me")
}

func (s *StorageMock) GetAll() ([]pool.Node, error) {
	panic("implement me")
}

func (s *StorageMock) Remove(pool.Node) error {
	panic("implement me")
}



func TestStrategy_Reserve_Positive(t *testing.T) {
	sm := new(StorageMock)
	expectedNode := *pool.NewNode(pool.NodeTypePersistent, "111", pool.NodeStatusBusy, "", time.Now().Unix(), 0, []pool.Capabilities{})
	sm.On("ReserveAvailable", mock.AnythingOfType("pool.Capabilities")).Return(expectedNode, nil)
	s := Strategy{sm}
	node, err := s.Reserve(pool.Capabilities{})
	assert.Equal(t, expectedNode, node)
	assert.Nil(t, err)
}

func TestStrategy_Reserve_Negative(t *testing.T) {
	sm := new(StorageMock)
	sm.On("ReserveAvailable", mock.AnythingOfType("pool.Capabilities")).Return(pool.Node{}, errors.New("Err"))
	s := Strategy{sm}
	_, err := s.Reserve(pool.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategy_CleanUp_Positive(t *testing.T) {
	sm := new(StorageMock)
	sm.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	s := Strategy{sm}
	node := *pool.NewNode(pool.NodeTypePersistent, "111", pool.NodeStatusBusy, "", time.Now().Unix(), 0, []pool.Capabilities{})
	err := s.CleanUp(node)
	assert.Nil(t, err)
}

func TestStrategy_CleanUp_Negative_NodeType(t *testing.T) {
	sm := new(StorageMock)
	sm.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	s := Strategy{sm}
	node := *pool.NewNode(pool.NodeTypeKubernetes, "111", pool.NodeStatusBusy, "", time.Now().Unix(), 0, []pool.Capabilities{})
	err := s.CleanUp(node)
	assert.NotNil(t, err)
}

func TestStrategy_CleanUp_Negative_NodeError(t *testing.T) {
	sm := new(StorageMock)
	sm.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(errors.New("Err"))
	s := Strategy{sm}
	node := *pool.NewNode(pool.NodeTypePersistent, "111", pool.NodeStatusBusy, "", time.Now().Unix(), 0, []pool.Capabilities{})
	err := s.CleanUp(node)
	assert.NotNil(t, err)
}
