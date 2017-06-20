package pool

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/qa-dev/jsonwire-grid/pool/strategy"
	"errors"
	"time"
)

type StrategyMock struct {
	mock.Mock
}

func (s *StrategyMock) Reserve(caps Capabilities) (Node, error) {
	args := s.Called(caps)
	return args.Get(0).(Node), args.Error(1)
}

func (s *StrategyMock) CleanUp(node Node) error {
	args := s.Called(node)
	return args.Error(0)
}

func TestNewStrategyList(t *testing.T) {
	a := assert.New(t)
	p := NewStrategyList([]StrategyInterface{})
	a.NotNil(new(StrategyInterface), p)
}

func TestStrategyList_Reserve_PositiveDirectOrder(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("Reserve", mock.AnythingOfType("pool.Capabilities")).Return(Node{}, strategy.ErrNotFound)
	s2 := new(StrategyListMock)
	expectedNode := *NewNode(NodeTypePersistent, "111", NodeStatusBusy, "", time.Now().Unix(), 0, []Capabilities{})
	s2.On("Reserve", mock.AnythingOfType("pool.Capabilities")).Return(expectedNode, nil)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	node, err := sl.Reserve(Capabilities{})
	assert.Nil(t, err)
	assert.Equal(t, expectedNode, node)
}

func TestStrategyList_Reserve_Positive_ReverseOrder(t *testing.T) {
	s1 := new(StrategyListMock)
	expectedNode := *NewNode(NodeTypePersistent, "111", NodeStatusBusy, "", time.Now().Unix(), 0, []Capabilities{})
	s1.On("Reserve", mock.AnythingOfType("pool.Capabilities")).Return(expectedNode, nil)
	s2 := new(StrategyListMock)
	s2.On("Reserve", mock.AnythingOfType("pool.Capabilities")).Return(Node{}, strategy.ErrNotFound)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	node, err := sl.Reserve(Capabilities{})
	assert.Nil(t, err)
	assert.Equal(t, expectedNode, node)
}

func TestStrategyList_Reserve_Negative_EmptyList(t *testing.T) {
	sl := NewStrategyList([]StrategyInterface{})
	_, err := sl.Reserve(Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategyList_Reserve_Negative_All_NotFound(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("Reserve", mock.AnythingOfType("pool.Capabilities")).Return(Node{}, strategy.ErrNotFound)
	s2 := new(StrategyListMock)
	s2.On("Reserve", mock.AnythingOfType("pool.Capabilities")).Return(Node{}, strategy.ErrNotFound)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	_, err := sl.Reserve(Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategyList_Reserve_Negative_Error(t *testing.T) {
	s1 := new(StrategyListMock)
	s1.On("Reserve", mock.AnythingOfType("pool.Capabilities")).Return(Node{}, errors.New("Error"))
	s2 := new(StrategyListMock)
	s2.On("Reserve", mock.AnythingOfType("pool.Capabilities")).Return(Node{}, nil)

	sl := NewStrategyList([]StrategyInterface{s1, s2})
	_, err := sl.Reserve(Capabilities{})
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
