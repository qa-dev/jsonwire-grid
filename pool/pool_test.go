package pool

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type StorageMock struct {
	mock.Mock
}

func (s *StorageMock) Add(node Node) error {
	args := s.Called(node)
	return args.Error(0)
}

func (s *StorageMock) ReserveAvailable(caps Capabilities) (Node, error) {
	args := s.Called(caps)
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

func TestNewPool(t *testing.T) {
	a := assert.New(t)
	p := NewPool(new(StorageMock))
	a.NotNil(new(Pool), p)
}

//---------------------------------

func TestPool_ReserveAvailableNode_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("ReserveAvailable", mock.AnythingOfType("pool.Capabilities")).Return(Node{}, nil)
	p := NewPool(s)
	node, err := p.ReserveAvailableNode(Capabilities{})
	a.NotNil(node)
	a.Nil(err)
}

func TestPool_ReserveAvailableNode_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("ReserveAvailable", mock.AnythingOfType("pool.Capabilities")).Return(Node{}, eError)
	p := NewPool(s)
	_, err := p.ReserveAvailableNode(Capabilities{})
	a.Error(err)
}

//---------------------------------

func TestPool_Add_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("Add", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(s)
	eAddress := "127.0.0.1"
	eNodeType := NodeTypeRegular
	err := p.Add(eNodeType, eAddress, []Capabilities{{"browserName": "ololo"}})
	a.Nil(err)
}

func TestPool_Add_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("Add", mock.AnythingOfType("pool.Node")).Return(eError)
	p := NewPool(s)
	eAddress := "127.0.0.1"
	eNodeType := NodeTypeRegular
	err := p.Add(eNodeType, eAddress, []Capabilities{})
	a.Error(err)
}

//---------------------------------

func TestPool_RegisterSession_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("SetBusy", mock.AnythingOfType("pool.Node"), mock.AnythingOfType("string")).Return(nil)
	p := NewPool(s)
	err := p.RegisterSession(new(Node), "testSessId")
	a.Nil(err)
}

func TestPool_RegisterSession_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("SetBusy", mock.AnythingOfType("pool.Node"), mock.AnythingOfType("string")).Return(eError)
	p := NewPool(s)
	err := p.RegisterSession(new(Node), "testSessId")
	a.Error(err)
	a.Equal(eError, err)
}

//---------------------------------

func TestPool_GetAll_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("GetAll").Return(make([]Node, 0), nil)
	p := NewPool(s)
	nodeList, err := p.GetAll()
	a.NotNil(nodeList)
	a.Nil(err)
}

func TestPool_GetAll_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("GetAll").Return(*new([]Node), eError)
	p := NewPool(s)
	nodeList, err := p.GetAll()
	a.Nil(nodeList)
	a.Error(err)
	a.Equal(eError, err)
}

//---------------------------------

func TestPool_GetNodeBySessionId_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("GetBySession", mock.AnythingOfType("string")).Return(Node{}, nil)
	p := NewPool(s)
	node, err := p.GetNodeBySessionId("testSessId")
	a.NotNil(node)
	a.Nil(err)
}

func TestPool_GetNodeBySessionId_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("GetBySession", mock.AnythingOfType("string")).Return(Node{}, eError)
	p := NewPool(s)
	node, err := p.GetNodeBySessionId("testSessId")
	a.Nil(node)
	a.Error(err)
}

//---------------------------------

func TestPool_GetNodeByAddress_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("GetByAddress", mock.AnythingOfType("string")).Return(Node{}, nil)
	p := NewPool(s)
	_, err := p.GetNodeByAddress("testAddress:testPort")
	a.Nil(err)
}

func TestPool_GetNodeByAddress_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("GetByAddress", mock.AnythingOfType("string")).Return(Node{}, eError)
	p := NewPool(s)
	_, err := p.GetNodeByAddress("testAddress:testPort")
	a.Error(err)
}

//---------------------------------

func TestPool_CleanUpNode_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(s)
	err := p.CleanUpNode(new(Node))
	a.Nil(err)
}

func TestPool_CleanUpNode_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(eError)
	p := NewPool(s)
	err := p.CleanUpNode(new(Node))
	a.Error(err)
}

//---------------------------------

func TestPool_Remove_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("Remove", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(s)
	err := p.Remove(new(Node))
	a.Nil(err)
}

func TestPool_Remove_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("Remove", mock.AnythingOfType("pool.Node")).Return(eError)
	p := NewPool(s)
	err := p.Remove(new(Node))
	a.Error(err)
}

//---------------------------------

func TestPool_CountNodes_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eCount := 0
	s.On("GetCountWithStatus", mock.AnythingOfType("*pool.NodeStatus")).Return(eCount, nil)
	p := NewPool(s)
	count, err := p.CountNodes(new(NodeStatus))
	a.Equal(eCount, count)
	a.Nil(err)
}

func TestPool_CountNodes_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("GetCountWithStatus", mock.AnythingOfType("*pool.NodeStatus")).Return(0, eError)
	p := NewPool(s)
	_, err := p.CountNodes(new(NodeStatus))
	a.Error(err)
}

//---------------------------------

func TestPool_fixNodeStatus_Positive_BusyExpired(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(s)
	node := NewNode(NodeTypeRegular, "", NodeStatusBusy, "", 0, 0, []Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.True(isFixed)
	a.Nil(err)
}

func TestPool_fixNodeStatus_Positive_ReservedExpired(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(s)
	node := NewNode(NodeTypeRegular, "", NodeStatusReserved, "", 0, 0, []Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.True(isFixed)
	a.Nil(err)
}

func TestPool_fixNodeStatus_Positive_BusyNotNotExpired(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(s)
	node := NewNode(NodeTypeRegular, "", NodeStatusBusy, "", time.Now().Unix(), 0, []Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.False(isFixed)
	a.Nil(err)
}

func TestPool_fixNodeStatus_Positive_ReservedNotNotExpired(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(s)
	node := NewNode(NodeTypeRegular, "", NodeStatusReserved, "", time.Now().Unix(), 0, []Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.False(isFixed)
	a.Nil(err)
}

func TestPool_fixNodeStatus_Positive_AvailableExpired(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(s)
	node := NewNode(NodeTypeRegular, "", NodeStatusAvailable, "", 0, 0, []Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.False(isFixed)
	a.Nil(err)
}

func TestPool_fixNodeStatus_NegativeBusy(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(eError)
	p := NewPool(s)
	node := NewNode(NodeTypeRegular, "", NodeStatusBusy, "", 0, 0, []Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.False(isFixed)
	a.Error(err)
}

func TestPool_fixNodeStatus_NegativeReserved(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(eError)
	p := NewPool(s)
	node := NewNode(NodeTypeRegular, "", NodeStatusReserved, "", 0, 0, []Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.False(isFixed)
	a.Error(err)
}

//---------------------------------
