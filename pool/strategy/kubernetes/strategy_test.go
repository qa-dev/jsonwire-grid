package kubernetes

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/qa-dev/jsonwire-grid/pool/strategy"
)

type providerMock struct {
	mock.Mock
}

func (p *providerMock) Create(podName string, nodeParams nodeParams) (nodeAddress string, err error) {
	args := p.Called(podName, nodeParams)
	return args.String(0), args.Error(1)
}

func (p *providerMock) Destroy(podName string) error {
	args := p.Called(podName)
	return args.Error(0)
}

func TestStrategy_Reserve_Positive(t *testing.T) {
	nodeCfg := nodeConfig{}
	nodeCfg.CapabilitiesList = []map[string]interface{}{{"cap1": "cal1"}}
	strategyConfig := strategyConfig{
		NodeList: []nodeConfig{nodeCfg},
	}
	sm := new(pool.StorageMock)
	sm.On("Add", mock.AnythingOfType("pool.Node"), mock.AnythingOfType("int")).Return(nil)
	sm.On("UpdateAddress", mock.AnythingOfType("pool.Node"), mock.AnythingOfType("string")).Return(nil)
	cm := new(capabilities.ComparatorMock)
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(true)
	pm := new(providerMock)
	expectedAddress := "addr"
	pm.On("Create", mock.AnythingOfType("string"), mock.AnythingOfType("nodeParams")).Return(expectedAddress, nil)
	str := Strategy{storage: sm, provider: pm, config: strategyConfig, capsComparator: cm}
	node, err := str.Reserve(capabilities.Capabilities{})
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, expectedAddress, node.Address)
}

func TestStrategy_Reserve_Negative_NotMatchCapabilities(t *testing.T) {
	eError := strategy.ErrNotFound
	s := Strategy{}
	_, err := s.Reserve(capabilities.Capabilities{})
	assert.Error(t, err, eError)
}

func TestStrategy_Reserve_Negative_ReserveAvailable(t *testing.T) {
	nodeCfg := nodeConfig{}
	nodeCfg.CapabilitiesList = []map[string]interface{}{{"cap1": "cal1"}}
	strategyConfig := strategyConfig{
		NodeList: []nodeConfig{nodeCfg},
	}
	sm := new(pool.StorageMock)
	sm.On("Add", mock.AnythingOfType("pool.Node"), mock.AnythingOfType("int")).Return(nil)
	cm := new(capabilities.ComparatorMock)
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(true)
	pm := new(providerMock)
	eError := errors.New("Error")
	pm.On("Create", mock.AnythingOfType("string"), mock.AnythingOfType("nodeParams")).Return("", eError)
	pm.On("Destroy", mock.AnythingOfType("string")).Return(nil)
	str := Strategy{storage: sm, provider: pm, config: strategyConfig, capsComparator: cm}
	_, err := str.Reserve(capabilities.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategy_Reserve_Negative_UpdateAddress(t *testing.T) {
	nodeCfg := nodeConfig{}
	nodeCfg.CapabilitiesList = []map[string]interface{}{{"cap1": "cal1"}}
	strategyConfig := strategyConfig{
		NodeList: []nodeConfig{nodeCfg},
	}
	sm := new(pool.StorageMock)
	sm.On("Add", mock.AnythingOfType("pool.Node"), mock.AnythingOfType("int")).Return(nil)
	sm.On("UpdateAddress", mock.AnythingOfType("pool.Node"), mock.AnythingOfType("string")).Return(errors.New("muhaha-error"))
	cm := new(capabilities.ComparatorMock)
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(true)
	pm := new(providerMock)
	eError := errors.New("Error")
	pm.On("Create", mock.AnythingOfType("string"), mock.AnythingOfType("nodeParams")).Return("", eError)
	pm.On("Destroy", mock.AnythingOfType("string")).Return(nil)
	str := Strategy{storage: sm, provider: pm, config: strategyConfig, capsComparator: cm}
	_, err := str.Reserve(capabilities.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategy_CleanUp_Positive(t *testing.T) {
	pm := new(providerMock)
	pm.On("Destroy", mock.AnythingOfType("string")).Return(nil)
	sm := new(pool.StorageMock)
	sm.On("Remove", mock.AnythingOfType("pool.Node")).Return(nil)
	s := Strategy{storage: sm, provider: pm}
	node := pool.Node{Key: "valid_key", Type: pool.NodeTypeKubernetes}
	err := s.CleanUp(node)
	assert.Nil(t, err)
}

func TestStrategy_CleanUp_Negative_NodeType(t *testing.T) {
	s := Strategy{}
	node := pool.Node{Type: pool.NodeTypePersistent}
	err := s.CleanUp(node)
	assert.Error(t, err, strategy.ErrNotApplicable)
}

func TestStrategy_CleanUp_Negative_EmptyNodeKey(t *testing.T) {
	s := Strategy{}
	node := pool.Node{Key: "", Type: pool.NodeTypeKubernetes} // empty node key
	err := s.CleanUp(node)
	assert.NotNil(t, err)
}

func TestStrategy_CleanUp_Negative_ErrProviderDestroy(t *testing.T) {
	pm := new(providerMock)
	pm.On("Destroy", mock.AnythingOfType("string")).Return(errors.New("Error"))
	s := Strategy{provider: pm}
	node := pool.Node{Type: pool.NodeTypeKubernetes, Address: "host:port"}
	err := s.CleanUp(node)
	assert.NotNil(t, err)
}

func TestStrategy_CleanUp_Negative_ErrStorageRemove(t *testing.T) {
	pm := new(providerMock)
	pm.On("Destroy", mock.AnythingOfType("string")).Return(nil)
	sm := new(pool.StorageMock)
	sm.On("Remove", mock.AnythingOfType("pool.Node")).Return(errors.New("Error"))
	s := Strategy{storage: sm, provider: pm}
	node := pool.Node{Type: pool.NodeTypeKubernetes, Address: "host:port"}
	err := s.CleanUp(node)
	assert.NotNil(t, err)
}

func TestStrategy_FixNodeStatus_Positive(t *testing.T) {
	pm := new(providerMock)
	pm.On("Destroy", mock.AnythingOfType("string")).Return(nil)
	sm := new(pool.StorageMock)
	sm.On("Remove", mock.AnythingOfType("pool.Node")).Return(nil)
	s := Strategy{storage: sm, provider: pm}
	node := pool.Node{Key: "valid_key", Type: pool.NodeTypeKubernetes}
	err := s.FixNodeStatus(node)
	assert.Nil(t, err)
}

func TestStrategy_FixNodeStatus_Negative_NodeType(t *testing.T) {
	s := Strategy{}
	node := pool.Node{Type: pool.NodeTypePersistent}
	err := s.FixNodeStatus(node)
	assert.Error(t, err, strategy.ErrNotApplicable)
}

func TestStrategy_FixNodeStatus_Negative_EmptyNodeKey(t *testing.T) {
	s := Strategy{}
	node := pool.Node{Key: "", Type: pool.NodeTypeKubernetes}
	err := s.FixNodeStatus(node)
	assert.NotNil(t, err)
}

func TestStrategy_FixNodeStatus_Negative_ErrProviderDestroy(t *testing.T) {
	pm := new(providerMock)
	pm.On("Destroy", mock.AnythingOfType("string")).Return(errors.New("Error"))
	s := Strategy{provider: pm}
	node := pool.Node{Type: pool.NodeTypeKubernetes, Address: "host:port"}
	err := s.FixNodeStatus(node)
	assert.NotNil(t, err)
}

func TestStrategy_FixNodeStatus_Negative_ErrStorageRemove(t *testing.T) {
	pm := new(providerMock)
	pm.On("Destroy", mock.AnythingOfType("string")).Return(nil)
	sm := new(pool.StorageMock)
	sm.On("Remove", mock.AnythingOfType("pool.Node")).Return(errors.New("Error"))
	s := Strategy{storage: sm, provider: pm}
	node := pool.Node{Type: pool.NodeTypeKubernetes, Address: "host:port"}
	err := s.FixNodeStatus(node)
	assert.NotNil(t, err)
}

func TestStrategy_findApplicableConfig_Positive(t *testing.T) {
	cm := new(capabilities.ComparatorMock)
	s := Strategy{capsComparator: cm}
	nodeCfg := nodeConfig{}
	nodeCfg.CapabilitiesList = []map[string]interface{}{{"cap1": "cal1"}}
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(true)
	applicableNodeCfg := s.findApplicableConfig([]nodeConfig{nodeCfg}, capabilities.Capabilities{"trololo": "lol"})
	assert.Equal(t, &nodeCfg, applicableNodeCfg)
}

func TestStrategy_findApplicableConfig_Negative(t *testing.T) {
	cm := new(capabilities.ComparatorMock)
	s := Strategy{capsComparator: cm}
	nodeCfg := nodeConfig{}
	nodeCfg.CapabilitiesList = []map[string]interface{}{{"cap1": "cal1"}}
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(false)
	applicableNodeCfg := s.findApplicableConfig([]nodeConfig{nodeCfg}, capabilities.Capabilities{"trololo": "lol"})
	assert.Nil(t, applicableNodeCfg)
}
