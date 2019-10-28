package persistent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/qa-dev/jsonwire-grid/pool/strategy"
)

type sessionsRemoverMockFactory struct {
	sessionsRemover sessionsRemover
}

func (f *sessionsRemoverMockFactory) create(abstractClient jsonwire.ClientInterface) sessionsRemover {
	return f.sessionsRemover
}

type sessionsRemoverMock struct {
	mock.Mock
}

func (r *sessionsRemoverMock) removeAllSessions() (int, error) {
	args := r.Called()
	return args.Int(0), args.Error(1)
}

func TestStrategy_Reserve_Positive(t *testing.T) {
	sm := new(pool.StorageMock)
	expectedNode := pool.Node{CapabilitiesList: []capabilities.Capabilities{{"cap1": "cal1"}}}
	sm.On("GetAll").Return([]pool.Node{expectedNode}, nil)
	cm := new(capabilities.ComparatorMock)
	cm.On("Register", mock.AnythingOfType("capabilities.Capabilities")).Return()
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(true)
	sm.On("ReserveAvailable", mock.AnythingOfType("[]pool.Node")).Return(expectedNode, nil)
	cfm := new(jsonwire.ClientFactoryMock)
	clm := new(jsonwire.ClientMock)
	cfm.On("Create", mock.AnythingOfType("string")).Return(clm)
	message := new(jsonwire.Message)
	clm.On("Health").Return(message, nil)
	srm := new(sessionsRemoverMock)
	srm.On("removeAllSessions").Return(0, nil)
	srfm := &sessionsRemoverMockFactory{srm}
	s := Strategy{storage: sm, capsComparator: cm, clientFactory: cfm, sessionsRemoverFactory: srfm}
	node, err := s.Reserve(capabilities.Capabilities{})
	assert.Nil(t, err)
	assert.Equal(t, expectedNode, node)
}

func TestStrategy_Reserve_Negative_GetAll_Error(t *testing.T) {
	sm := new(pool.StorageMock)
	eError := errors.New("Error")
	sm.On("GetAll").Return([]pool.Node{}, eError)
	s := Strategy{storage: sm}
	_, err := s.Reserve(capabilities.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategy_Reserve_Negative_NotMatchCapabilities(t *testing.T) {
	sm := new(pool.StorageMock)
	eError := strategy.ErrNotFound
	sm.On("GetAll").Return([]pool.Node{{}}, nil) // >= 1
	s := Strategy{storage: sm}
	_, err := s.Reserve(capabilities.Capabilities{})
	assert.Error(t, err, eError)
}

func TestStrategy_Reserve_Negative_ReserveAvailable(t *testing.T) {
	sm := new(pool.StorageMock)
	eError := errors.New("Error")
	sm.On("GetAll").Return([]pool.Node{{CapabilitiesList: []capabilities.Capabilities{{}}}}, nil)
	cm := new(capabilities.ComparatorMock)
	cm.On("Register", mock.AnythingOfType("capabilities.Capabilities")).Return()
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(true)
	sm.On("ReserveAvailable", mock.AnythingOfType("[]pool.Node")).Return(pool.Node{}, eError)
	s := Strategy{storage: sm, capsComparator: cm}
	_, err := s.Reserve(capabilities.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategy_Reserve_Negative_Client_Status_Error(t *testing.T) {
	sm := new(pool.StorageMock)
	sm.On("GetAll").Return([]pool.Node{{}}, nil)
	cm := new(capabilities.ComparatorMock)
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(true)
	expectedNode := pool.Node{CapabilitiesList: []capabilities.Capabilities{{"cap1": "cal1"}}}
	sm.On("ReserveAvailable", mock.AnythingOfType("[]pool.Node")).Return([]pool.Node{expectedNode}, nil)
	cfm := new(jsonwire.ClientFactoryMock)
	clm := new(jsonwire.ClientMock)
	cfm.On("Create", mock.AnythingOfType("string")).Return(clm)
	eError := errors.New("Error")
	clm.On("Health").Return(nil, eError)
	s := Strategy{storage: sm, capsComparator: cm, clientFactory: cfm}
	_, err := s.Reserve(capabilities.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategy_Reserve_Negative_Client_Status_NotOk(t *testing.T) {
	sm := new(pool.StorageMock)
	sm.On("GetAll").Return([]pool.Node{{}}, nil)
	cm := new(capabilities.ComparatorMock)
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(true)
	expectedNode := pool.Node{CapabilitiesList: []capabilities.Capabilities{{"cap1": "cal1"}}}
	sm.On("ReserveAvailable", mock.AnythingOfType("[]pool.Node")).Return([]pool.Node{expectedNode}, nil)
	cfm := new(jsonwire.ClientFactoryMock)
	clm := new(jsonwire.ClientMock)
	cfm.On("Create", mock.AnythingOfType("string")).Return(clm)
	message := new(jsonwire.Message)
	message.Status = -99
	clm.On("Health").Return(message, nil)
	s := Strategy{storage: sm, capsComparator: cm, clientFactory: cfm}
	_, err := s.Reserve(capabilities.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategy_Reserve_Negative_removeAllSessions_Error(t *testing.T) {
	sm := new(pool.StorageMock)
	sm.On("GetAll").Return([]pool.Node{{}}, nil)
	cm := new(capabilities.ComparatorMock)
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(true)
	expectedNode := pool.Node{CapabilitiesList: []capabilities.Capabilities{{"cap1": "cal1"}}}
	sm.On("ReserveAvailable", mock.AnythingOfType("[]pool.Node")).Return([]pool.Node{expectedNode}, nil)
	cfm := new(jsonwire.ClientFactoryMock)
	clm := new(jsonwire.ClientMock)
	cfm.On("Create", mock.AnythingOfType("string")).Return(clm)
	message := new(jsonwire.Message)
	clm.On("Health").Return(message, nil)
	eError := errors.New("Error")
	srm := new(sessionsRemoverMock)
	srm.On("removeAllSessions").Return(0, eError)
	srfm := &sessionsRemoverMockFactory{srm}
	s := Strategy{storage: sm, capsComparator: cm, clientFactory: cfm, sessionsRemoverFactory: srfm}
	_, err := s.Reserve(capabilities.Capabilities{})
	assert.NotNil(t, err)
}

func TestStrategy_CleanUp_Positive(t *testing.T) {
	sm := new(pool.StorageMock)
	sm.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	s := Strategy{storage: sm}
	node := pool.Node{Type: pool.NodeTypePersistent}
	err := s.CleanUp(node)
	assert.Nil(t, err)
}

func TestStrategy_CleanUp_Negative_NodeType(t *testing.T) {
	s := Strategy{}
	node := pool.Node{Type: pool.NodeTypeKubernetes}
	err := s.CleanUp(node)
	assert.Error(t, err, strategy.ErrNotApplicable)
}

func TestStrategy_CleanUp_Negative_NodeError(t *testing.T) {
	sm := new(pool.StorageMock)
	sm.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(errors.New("Err"))
	s := Strategy{storage: sm}
	node := pool.Node{Type: pool.NodeTypePersistent}
	err := s.CleanUp(node)
	assert.NotNil(t, err)
}

func TestStrategy_FixNodeStatus_Positive(t *testing.T) {
	sm := new(pool.StorageMock)
	sm.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	s := Strategy{storage: sm}
	node := pool.Node{Type: pool.NodeTypePersistent}
	err := s.FixNodeStatus(node)
	assert.Nil(t, err)
}

func TestStrategy_FixNodeStatus_Negative_NodeType(t *testing.T) {
	sm := new(pool.StorageMock)
	sm.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(nil)
	s := Strategy{storage: sm}
	node := pool.Node{Type: pool.NodeTypeKubernetes}
	err := s.FixNodeStatus(node)
	assert.NotNil(t, err)
}

func TestStrategy_FixNodeStatus_Negative_NodeError(t *testing.T) {
	sm := new(pool.StorageMock)
	sm.On("SetAvailable", mock.AnythingOfType("pool.Node")).Return(errors.New("Err"))
	s := Strategy{storage: sm}
	node := pool.Node{Type: pool.NodeTypePersistent}
	err := s.FixNodeStatus(node)
	assert.NotNil(t, err)
}

func TestStrategy_findApplicableNodes_Positive(t *testing.T) {
	cm := new(capabilities.ComparatorMock)
	s := Strategy{capsComparator: cm}
	caps := []capabilities.Capabilities{{"caps1": "val1"}}
	nodeList := []pool.Node{
		{Address: "node1", CapabilitiesList: caps},
		{Address: "node2", CapabilitiesList: caps},
	}
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(true)
	applicableList := s.findApplicableNodes(nodeList, capabilities.Capabilities{"trololo": "lol"})
	assert.Len(t, applicableList, 2)
}

func TestStrategy_findApplicableNodes_Negative(t *testing.T) {
	cm := new(capabilities.ComparatorMock)
	s := Strategy{capsComparator: cm}
	caps := []capabilities.Capabilities{{"caps1": "val1"}}
	nodeList := []pool.Node{
		{Address: "node1", CapabilitiesList: caps},
	}
	cm.On("Compare", mock.AnythingOfType("capabilities.Capabilities"), mock.AnythingOfType("capabilities.Capabilities")).Return(false)
	applicableList := s.findApplicableNodes(nodeList, capabilities.Capabilities{"trololo": "lol"})
	assert.Len(t, applicableList, 0)
}
