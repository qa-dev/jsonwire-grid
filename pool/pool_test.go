package pool

import (
	"errors"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/qa-dev/jsonwire-grid/pool/strategy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	a := assert.New(t)
	p := NewPool(new(StorageMock), new(StrategyListMock))
	a.NotNil(new(Pool), p)
}

//---------------------------------

func TestPool_ReserveAvailableNode_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StrategyListMock)
	s.On("Reserve", mock.AnythingOfType("capabilities.Capabilities")).Return(Node{}, nil)
	p := NewPool(new(StorageMock), s)
	node, err := p.ReserveAvailableNode(capabilities.Capabilities{})
	a.NotNil(node)
	a.Nil(err)
}

func TestPool_ReserveAvailableNode_Negative(t *testing.T) {
	a := assert.New(t)
	eError := strategy.ErrNotFound
	s := new(StrategyListMock)
	s.On("Reserve", mock.AnythingOfType("capabilities.Capabilities")).Return(Node{}, eError)
	p := NewPool(new(StorageMock), s)
	_, err := p.ReserveAvailableNode(capabilities.Capabilities{})
	a.Error(err)
}

//---------------------------------

func TestPool_Add_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("Add", mock.AnythingOfType("pool.Node"), mock.AnythingOfType("int")).Return(nil)
	p := NewPool(s, new(StrategyListMock))
	eAddress := "127.0.0.1"
	eNodeType := NodeTypePersistent
	err := p.Add(eNodeType, eAddress, []capabilities.Capabilities{{"browserName": "ololo"}})
	a.Nil(err)
}

func TestPool_Add_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("Add", mock.AnythingOfType("pool.Node")).Return(eError)
	p := NewPool(s, new(StrategyListMock))
	eAddress := "127.0.0.1"
	eNodeType := NodeTypePersistent
	err := p.Add(eNodeType, eAddress, []capabilities.Capabilities{})
	a.Error(err)
}

//---------------------------------

func TestPool_RegisterSession_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("SetBusy", mock.AnythingOfType("pool.Node"), mock.AnythingOfType("string")).Return(nil)
	p := NewPool(s, new(StrategyListMock))
	err := p.RegisterSession(new(Node), "testSessId")
	a.Nil(err)
}

func TestPool_RegisterSession_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("SetBusy", mock.AnythingOfType("pool.Node"), mock.AnythingOfType("string")).Return(eError)
	p := NewPool(s, new(StrategyListMock))
	err := p.RegisterSession(new(Node), "testSessId")
	a.Error(err)
	a.Equal(eError, err)
}

//---------------------------------

func TestPool_GetAll_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("GetAll").Return(make([]Node, 0), nil)
	p := NewPool(s, new(StrategyListMock))
	nodeList, err := p.GetAll()
	a.NotNil(nodeList)
	a.Nil(err)
}

func TestPool_GetAll_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("GetAll").Return(*new([]Node), eError)
	p := NewPool(s, new(StrategyListMock))
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
	p := NewPool(s, new(StrategyListMock))
	node, err := p.GetNodeBySessionID("testSessId")
	a.NotNil(node)
	a.Nil(err)
}

func TestPool_GetNodeBySessionId_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("GetBySession", mock.AnythingOfType("string")).Return(Node{}, eError)
	p := NewPool(s, new(StrategyListMock))
	node, err := p.GetNodeBySessionID("testSessId")
	a.Nil(node)
	a.Error(err)
}

//---------------------------------

func TestPool_GetNodeByAddress_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("GetByAddress", mock.AnythingOfType("string")).Return(Node{}, nil)
	p := NewPool(s, new(StrategyListMock))
	_, err := p.GetNodeByAddress("testAddress:testPort")
	a.Nil(err)
}

func TestPool_GetNodeByAddress_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("GetByAddress", mock.AnythingOfType("string")).Return(Node{}, eError)
	p := NewPool(s, new(StrategyListMock))
	_, err := p.GetNodeByAddress("testAddress:testPort")
	a.Error(err)
}

//---------------------------------

func TestPool_CleanUpNode_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StrategyListMock)
	s.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(new(StorageMock), s)
	err := p.CleanUpNode(new(Node))
	a.Nil(err)
}

func TestPool_CleanUpNode_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StrategyListMock)
	eError := errors.New("Error")
	s.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(eError)
	p := NewPool(new(StorageMock), s)
	err := p.CleanUpNode(new(Node))
	a.Error(err)
}

//---------------------------------

func TestPool_Remove_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	s.On("Remove", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(s, new(StrategyListMock))
	err := p.Remove(new(Node))
	a.Nil(err)
}

func TestPool_Remove_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("Remove", mock.AnythingOfType("pool.Node")).Return(eError)
	p := NewPool(s, new(StrategyListMock))
	err := p.Remove(new(Node))
	a.Error(err)
}

//---------------------------------

func TestPool_CountNodes_Positive(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eCount := 0
	s.On("GetCountWithStatus", mock.AnythingOfType("*pool.NodeStatus")).Return(eCount, nil)
	p := NewPool(s, new(StrategyListMock))
	count, err := p.CountNodes(new(NodeStatus))
	a.Equal(eCount, count)
	a.Nil(err)
}

func TestPool_CountNodes_Negative(t *testing.T) {
	a := assert.New(t)
	s := new(StorageMock)
	eError := errors.New("Error")
	s.On("GetCountWithStatus", mock.AnythingOfType("*pool.NodeStatus")).Return(0, eError)
	p := NewPool(s, new(StrategyListMock))
	_, err := p.CountNodes(new(NodeStatus))
	a.Error(err)
}

//---------------------------------

func TestPool_fixNodeStatus_Positive_BusyExpired(t *testing.T) {
	a := assert.New(t)
	slm := new(StrategyListMock)
	slm.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(new(StorageMock), slm)
	node := NewNode(NodeTypePersistent, "", NodeStatusBusy, "", 0, 0, []capabilities.Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.True(isFixed)
	a.Nil(err)
}

func TestPool_fixNodeStatus_Positive_ReservedExpired(t *testing.T) {
	a := assert.New(t)
	slm := new(StrategyListMock)
	slm.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(new(StorageMock), slm)
	node := NewNode(NodeTypePersistent, "", NodeStatusReserved, "", 0, 0, []capabilities.Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.True(isFixed)
	a.Nil(err)
}

func TestPool_fixNodeStatus_Positive_BusyNotNotExpired(t *testing.T) {
	a := assert.New(t)
	slm := new(StrategyListMock)
	slm.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(new(StorageMock), slm)
	node := NewNode(NodeTypePersistent, "", NodeStatusBusy, "", time.Now().Unix(), 0, []capabilities.Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.False(isFixed)
	a.Nil(err)
}

func TestPool_fixNodeStatus_Positive_ReservedNotNotExpired(t *testing.T) {
	a := assert.New(t)
	slm := new(StrategyListMock)
	slm.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(new(StorageMock), slm)
	node := NewNode(NodeTypePersistent, "", NodeStatusReserved, "", time.Now().Unix(), 0, []capabilities.Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.False(isFixed)
	a.Nil(err)
}

func TestPool_fixNodeStatus_Positive_AvailableExpired(t *testing.T) {
	a := assert.New(t)
	slm := new(StrategyListMock)
	slm.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(nil)
	p := NewPool(new(StorageMock), slm)
	node := NewNode(NodeTypePersistent, "", NodeStatusAvailable, "", 0, 0, []capabilities.Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.False(isFixed)
	a.Nil(err)
}

func TestPool_fixNodeStatus_NegativeBusy(t *testing.T) {
	a := assert.New(t)
	eError := errors.New("Error")
	slm := new(StrategyListMock)
	slm.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(eError)
	p := NewPool(new(StorageMock), slm)
	node := NewNode(NodeTypePersistent, "", NodeStatusBusy, "", 0, 0, []capabilities.Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.False(isFixed)
	a.Error(err)
}

func TestPool_fixNodeStatus_NegativeReserved(t *testing.T) {
	a := assert.New(t)
	eError := errors.New("Error")
	slm := new(StrategyListMock)
	slm.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(eError)
	p := NewPool(new(StorageMock), slm)
	node := NewNode(NodeTypePersistent, "", NodeStatusReserved, "", 0, 0, []capabilities.Capabilities{})
	isFixed, err := p.fixNodeStatus(node)
	a.False(isFixed)
	a.Error(err)
}

//---------------------------------
