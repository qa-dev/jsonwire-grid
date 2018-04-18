package local

import (
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage_Add_Positive(t *testing.T) {
	s := Storage{db: map[string]*pool.Node{}}
	err := s.Add(pool.Node{}, 0)
	assert.NoError(t, err)
	assert.Len(t, s.db, 1)
}

func TestStorage_Add_Positive_Repeat(t *testing.T) {
	s := Storage{db: map[string]*pool.Node{"1": {Key: "1"}}}
	err := s.Add(pool.Node{Key: "1"}, 0)
	assert.NoError(t, err)
	assert.Len(t, s.db, 1)
}

func TestStorage_Add_Negative_LimitReached(t *testing.T) {
	s := Storage{db: map[string]*pool.Node{"1": {Key: "1"}}}
	limit := 1
	err := s.Add(pool.Node{Key: "2"}, limit)
	assert.Error(t, err, "limit reached")
	assert.Len(t, s.db, limit)
}

func TestStorage_ReserveAvailable_Positive(t *testing.T) {
	expectedNode := pool.Node{Key: "1", Status: pool.NodeStatusAvailable}
	s := Storage{db: map[string]*pool.Node{expectedNode.Key: &expectedNode}}
	node, err := s.ReserveAvailable([]pool.Node{expectedNode})
	assert.NoError(t, err)
	assert.Equal(t, expectedNode, node)
	assert.Equal(t, pool.NodeStatusReserved, s.db[node.Key].Status)
}

func TestStorage_ReserveAvailable_Negative_NotFoundAvailableNodes(t *testing.T) {
	expectedNode := pool.Node{Key: "1", Status: pool.NodeStatusBusy}
	s := Storage{db: map[string]*pool.Node{expectedNode.Key: &expectedNode}}
	_, err := s.ReserveAvailable([]pool.Node{expectedNode})
	assert.Error(t, err, storage.ErrNotFound)
}

func TestStorage_ReserveAvailable_Negative_InvalidNodeList(t *testing.T) {
	s := Storage{db: map[string]*pool.Node{}}
	_, err := s.ReserveAvailable([]pool.Node{{Key: "awd"}})
	assert.Error(t, err, storage.ErrNotFound)
}

func TestStorage_SetBusy_Positive(t *testing.T) {
	expectedNode := pool.Node{Key: "1"}
	s := Storage{db: map[string]*pool.Node{expectedNode.Key: &expectedNode}}
	expectedSessionID := "expectedSessionID"
	err := s.SetBusy(expectedNode, expectedSessionID)
	assert.NoError(t, err)
	assert.Equal(t, pool.NodeStatusBusy, s.db[expectedNode.Key].Status)
	assert.Equal(t, expectedSessionID, s.db[expectedNode.Key].SessionID)
}

func TestStorage_SetBusy_Negative(t *testing.T) {
	expectedNode := pool.Node{Key: "1"}
	s := Storage{db: map[string]*pool.Node{}}
	expectedSessionID := "expectedSessionID"
	err := s.SetBusy(expectedNode, expectedSessionID)
	assert.Error(t, err, storage.ErrNotFound)
}

func TestStorage_SetAvailable_Positive(t *testing.T) {
	expectedNode := pool.Node{Key: "1"}
	s := Storage{db: map[string]*pool.Node{expectedNode.Key: &expectedNode}}
	err := s.SetAvailable(expectedNode)
	assert.NoError(t, err)
	assert.Equal(t, pool.NodeStatusAvailable, s.db[expectedNode.Key].Status)
}

func TestStorage_SetAvailable_Negative(t *testing.T) {
	expectedNode := pool.Node{Key: "1"}
	s := Storage{db: map[string]*pool.Node{}}
	err := s.SetAvailable(expectedNode)
	assert.Error(t, err, storage.ErrNotFound)
}

func TestStorage_GetCountWithStatus_Positive_All(t *testing.T) {
	s := Storage{db: map[string]*pool.Node{"1": {Key: "1"}, "2": {Key: "2"}}}
	count, err := s.GetCountWithStatus(nil)
	assert.NoError(t, err)
	assert.Equal(t, count, len(s.db))
}

func TestStorage_GetCountWithStatus_Positive_One(t *testing.T) {
	expectedStatus := pool.NodeStatusBusy
	s := Storage{db: map[string]*pool.Node{"1": {Key: "1", Status: expectedStatus}, "2": {Key: "2"}}}
	count, err := s.GetCountWithStatus(&expectedStatus)
	assert.NoError(t, err)
	assert.Equal(t, count, 1)
}

func TestStorage_GetBySession_Positive(t *testing.T) {
	expectedNode := pool.Node{Key: "1"}
	s := Storage{db: map[string]*pool.Node{expectedNode.Key: &expectedNode}}
	node, err := s.GetBySession(expectedNode.SessionID)
	assert.NoError(t, err)
	assert.Equal(t, expectedNode, node)
}

func TestStorage_GetBySession_Negative(t *testing.T) {
	expectedNode := pool.Node{Key: "1"}
	s := Storage{db: map[string]*pool.Node{}}
	_, err := s.GetBySession(expectedNode.SessionID)
	assert.Error(t, err, storage.ErrNotFound)
}

func TestStorage_GetByAddress_Positive(t *testing.T) {
	expectedNode := pool.Node{Address: "1", Key: "12d"}
	s := Storage{db: map[string]*pool.Node{expectedNode.Key: &expectedNode}}
	node, err := s.GetByAddress(expectedNode.Address)
	assert.NoError(t, err)
	assert.Equal(t, expectedNode, node)
}

func TestStorage_GetByAddress_Negative(t *testing.T) {
	expectedNode := pool.Node{Address: "1", Key: "1234567890"}
	s := Storage{db: map[string]*pool.Node{}}
	_, err := s.GetByAddress(expectedNode.Address)
	assert.Error(t, err, storage.ErrNotFound)
}

func TestStorage_GetAll_Positive(t *testing.T) {
	s := Storage{db: map[string]*pool.Node{"1": {Key: "1"}, "2": {Key: "2"}}}
	nodeList, err := s.GetAll()
	assert.NoError(t, err)
	assert.Len(t, nodeList, 2)
}

func TestStorage_Remove_Positive(t *testing.T) {
	node := pool.Node{Key: "1"}
	s := Storage{db: map[string]*pool.Node{node.Key: &node}}
	err := s.Remove(node)
	assert.NoError(t, err)
}

func TestStorage_Remove_Negative(t *testing.T) {
	s := Storage{db: map[string]*pool.Node{}}
	node := pool.Node{Key: "1"}
	err := s.Remove(node)
	assert.Error(t, err, storage.ErrNotFound)
}

func TestStorage_UpdateAddress_UpdatesValue(t *testing.T) {
	key := "1234567890"
	node := pool.Node{Address: "1", Key: key}
	s := Storage{db: map[string]*pool.Node{key: &node}}
	expectedAddress := "2"
	err := s.UpdateAddress(node, expectedAddress)
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, s.db[key].Address)
}

func TestStorage_UpdateAddress_ReturnsErrNotFound(t *testing.T) {
	key := "1234567890"
	node := pool.Node{Address: "1", Key: key}
	s := Storage{db: map[string]*pool.Node{key: &node}}
	err := s.UpdateAddress(pool.Node{Key:"12345"}, "1234567890")
	assert.Equal(t, storage.ErrNotFound, err)
}