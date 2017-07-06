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

func TestNewStrategyList(t *testing.T) {
	a := assert.New(t)
	p := NewStrategyList([]StrategyInterface{})
	a.NotNil(new(StrategyInterface), p)
}

func TestStrategyList_Reserve_PositiveDirectOrder(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("Reserve", mock.AnythingOfType("capabilities.Capabilities")).Return(Node{}, strategy.ErrNotFound)
	s2 := new(StrategyListMock)
	expectedNode := *NewNode(NodeTypePersistent, "111", NodeStatusBusy, "", time.Now().Unix(), 0, []capabilities.Capabilities{})
	s2.On("Reserve", mock.AnythingOfType("capabilities.Capabilities")).Return(expectedNode, nil)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	node, err := sl.Reserve(capabilities.Capabilities{})
	assert.Nil(t, err)
	assert.Equal(t, expectedNode, node)
}

func TestStrategyList_Reserve_Positive_ReverseOrder(t *testing.T) {
	s1 := new(StrategyListMock)
	expectedNode := *NewNode(NodeTypePersistent, "111", NodeStatusBusy, "", time.Now().Unix(), 0, []capabilities.Capabilities{})
	s1.On("Reserve", mock.AnythingOfType("capabilities.Capabilities")).Return(expectedNode, nil)
	s2 := new(StrategyListMock)
	s2.On("Reserve", mock.AnythingOfType("capabilities.Capabilities")).Return(Node{}, strategy.ErrNotFound)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	node, err := sl.Reserve(capabilities.Capabilities{})
	assert.Nil(t, err)
	assert.Equal(t, expectedNode, node)
}

func TestStrategyList_Reserve_Negative_EmptyList(t *testing.T) {
	sl := NewStrategyList([]StrategyInterface{})
	_, err := sl.Reserve(capabilities.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategyList_Reserve_Negative_All_NotFound(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("Reserve", mock.AnythingOfType("capabilities.Capabilities")).Return(Node{}, strategy.ErrNotFound)
	s2 := new(StrategyListMock)
	s2.On("Reserve", mock.AnythingOfType("capabilities.Capabilities")).Return(Node{}, strategy.ErrNotFound)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	_, err := sl.Reserve(capabilities.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategyList_Reserve_Negative_Error(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("Reserve", mock.AnythingOfType("capabilities.Capabilities")).Return(Node{}, errors.New("Error"))
	s2 := new(StrategyListMock)
	s2.On("Reserve", mock.AnythingOfType("capabilities.Capabilities")).Return(Node{}, nil)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	_, err := sl.Reserve(capabilities.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategyList_CleanUp_Positive_DirectOrder(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(strategy.ErrNotApplicable)
	s2 := new(StrategyListMock)
	s2.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(nil)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	err := sl.CleanUp(Node{})
	assert.Nil(t, err)
}

func TestStrategyList_CleanUp_Positive_ReverseOrder(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(nil)
	s2 := new(StrategyListMock)
	s2.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(strategy.ErrNotApplicable)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	err := sl.CleanUp(Node{})
	assert.Nil(t, err)
}

func TestStrategyList_CleanUp_Negative_EmptyList(t *testing.T) {
	sl := NewStrategyList([]StrategyInterface{})
	err := sl.CleanUp(Node{})
	assert.NotNil(t, err)
}

func TestStrategyList_CleanUp_Negative_All_NotApplicable(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(strategy.ErrNotApplicable)
	s2 := new(StrategyListMock)
	s2.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(strategy.ErrNotApplicable)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	err := sl.CleanUp(Node{})
	assert.NotNil(t, err)
}

func TestStrategyList_CleanUp_Negative_Error(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(errors.New("Error"))
	s2 := new(StrategyListMock)
	s2.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(nil)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	err := sl.CleanUp(Node{})
	assert.NotNil(t, err)
}

func TestStrategyList_FixNodeStatus_Positive_DirectOrder(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(strategy.ErrNotApplicable)
	s2 := new(StrategyListMock)
	s2.On("CleanUp", mock.AnythingOfType("pool.Node")).Return(nil)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	err := sl.CleanUp(Node{})
	assert.Nil(t, err)
}

func TestStrategyList_FixNodeStatus_Positive_ReverseOrder(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(nil)
	s2 := new(StrategyListMock)
	s2.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(strategy.ErrNotApplicable)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	err := sl.FixNodeStatus(Node{})
	assert.Nil(t, err)
}

func TestStrategyList_FixNodeStatus_Negative_EmptyList(t *testing.T) {
	sl := NewStrategyList([]StrategyInterface{})
	err := sl.FixNodeStatus(Node{})
	assert.NotNil(t, err)
}

func TestStrategyList_FixNodeStatus_Negative_All_NotApplicable(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(strategy.ErrNotApplicable)
	s2 := new(StrategyListMock)
	s2.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(strategy.ErrNotApplicable)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	err := sl.FixNodeStatus(Node{})
	assert.NotNil(t, err)
}

func TestStrategyList_FixNodeStatus_Negative_Error(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(errors.New("Error"))
	s2 := new(StrategyListMock)
	s2.On("FixNodeStatus", mock.AnythingOfType("pool.Node")).Return(nil)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	err := sl.FixNodeStatus(Node{})
	assert.NotNil(t, err)
}
