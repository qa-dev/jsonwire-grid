package tests

import (
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

// Сделано для того чтобы тестировать одним набором тестов все реализации
type PrepareInterface interface {
	SetUp()
	TearDown()
	CreateStorage() (pool.StorageInterface, func())
}

// testStorage_Add проверка корректости добавления ноды в хранилище
func testStorage_Add(t *testing.T, p PrepareInterface) {
	t.Parallel()

	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	expectedNode := pool.Node{
		Key: "111",
		Address: "address1",
		CapabilitiesList: []capabilities.Capabilities{
			{"trololo": "lolo"},
			{
				"mysupercapability1": "mysupercapabilityValue1",
				"mysupercapability2": "mysupercapabilityValue2",
			},
		},
	}
	err := storage.Add(expectedNode, 0)
	assert.Nil(t, err)
	nodeList, err := storage.GetAll()
	if err != nil {
		t.Fatal("Error get all nodes list, " + err.Error())
	}
	assert.Len(t, nodeList, 1, "Added more than one node")
	assert.Equal(t, expectedNode.Key, nodeList[0].Key)
	assert.Equal(t, expectedNode.Type, nodeList[0].Type)
	assert.Equal(t, expectedNode.Address, nodeList[0].Address)
	assert.Equal(t, expectedNode.Status, nodeList[0].Status)
	assert.Equal(t, expectedNode.SessionID, nodeList[0].SessionID)
	assert.Equal(t, expectedNode.Updated, nodeList[0].Updated)
	assert.Equal(t, expectedNode.Registered, nodeList[0].Registered)
	//assert.Equal(t, expectedNode.CapabilitiesList, nodeList[0].CapabilitiesList) //todo: доделать
}

// testStorage_Add_Repeat проверка того что при повторном добавлении ноды вместо дублирования происходит корректный апдейт
func testStorage_Add_Repeat(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	node := pool.Node{
		Address:          "ololo",
		CapabilitiesList: []capabilities.Capabilities{{"trololo": "lolo"}},
	}
	err := storage.Add(node, 0)
	assert.Nil(t, err)
	node.SessionID = "changedSessionID"
	err = storage.Add(node, 0)
	assert.Nil(t, err)
	nodeList, err := storage.GetAll()
	assert.Nil(t, err)
	assert.Len(t, nodeList, 1, "Added more than one node")
	assert.Equal(t, node.SessionID, nodeList[0].SessionID, "SessionId not updated")
	//todo: доделать capabilities
}

// testStorage_Add_Limit_Overflow проверка того что при переполнении лимита, запись не добавляется в хранилище
func testStorage_Add_Limit_Overflow(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	node := pool.Node{
		Address:          "ololo",
		CapabilitiesList: []capabilities.Capabilities{{"trololo": "lolo"}},
		Type:             pool.NodeTypePersistent,
	}
	limit := 1
	err := storage.Add(node, limit)
	assert.Nil(t, err)
	node.Address = "ololo1"
	err = storage.Add(node, limit)
	assert.NotNil(t, err)
	nodeList, err := storage.GetAll()
	assert.Nil(t, err)
	assert.Len(t, nodeList, limit, "Added more than "+strconv.Itoa(limit)+"one node")
}

// testStorage_GetAll проверка получения всех нод
func testStorage_GetAll(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	expectedNodeList := make([]pool.Node, 0)
	for _, addr := range []string{"addr1", "addr2"} {
		node := pool.Node{
			Key: addr,
			Address:          addr,
			CapabilitiesList: []capabilities.Capabilities{{"trololo": "lolo"}},
		}
		expectedNodeList = append(expectedNodeList, node)
		err := storage.Add(node, 0)
		if err != nil {
			t.Fatal("Error add node, " + err.Error())
		}
	}
	nodeList, err := storage.GetAll()
	assert.Nil(t, err)
	assert.Len(t, nodeList, len(expectedNodeList))
	for _, expectedNode := range expectedNodeList {
		isNodeMatch := false
		for _, node := range nodeList {
			if node.Address == expectedNode.Address {
				assert.Equal(t, expectedNode.Key, node.Key)
				assert.Equal(t, expectedNode.Type, node.Type)
				assert.Equal(t, expectedNode.Address, node.Address)
				assert.Equal(t, expectedNode.Status, node.Status)
				assert.Equal(t, expectedNode.SessionID, node.SessionID)
				assert.Equal(t, expectedNode.Updated, node.Updated)
				assert.Equal(t, expectedNode.Registered, node.Registered)
				//todo: доделать capabilities
				isNodeMatch = true
			}
		}
		assert.True(t, isNodeMatch, "Not expected node not found in nodes list")

	}
}

// testStorage_GetByAddress проверка получения ноды по адресу
func testStorage_GetByAddress(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	expectedNode := pool.Node{Address: "mySuperAddress"}
	err := storage.Add(expectedNode, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	node, err := storage.GetByAddress(expectedNode.Address)
	assert.Nil(t, err)
	assert.Equal(t, expectedNode.Type, node.Type)
	assert.Equal(t, expectedNode.Address, node.Address)
	assert.Equal(t, expectedNode.Status, node.Status)
	assert.Equal(t, expectedNode.SessionID, node.SessionID)
	assert.Equal(t, expectedNode.Updated, node.Updated)
	assert.Equal(t, expectedNode.Registered, node.Registered)

}

// testStorage_GetBySession проверка получения ноды по sessionId
func testStorage_GetBySession(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	expectedNode := pool.Node{Address: "mySuperAddress"}
	err := storage.Add(expectedNode, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	node, err := storage.GetBySession(expectedNode.SessionID)
	assert.Nil(t, err)
	assert.Equal(t, expectedNode.Type, node.Type)
	assert.Equal(t, expectedNode.Address, node.Address)
	assert.Equal(t, expectedNode.Status, node.Status)
	assert.Equal(t, expectedNode.SessionID, node.SessionID)
	assert.Equal(t, expectedNode.Updated, node.Updated)
	assert.Equal(t, expectedNode.Registered, node.Registered)

}

// testStorage_GetCountWithStatus проверка получения колличества нод с определенным статусом
func testStorage_GetCountWithStatus(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	err := storage.Add(pool.Node{Status: pool.NodeStatusAvailable, Address: "1"}, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	err = storage.Add(pool.Node{Status: pool.NodeStatusAvailable, Address: "2"}, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	err = storage.Add(pool.Node{Status: pool.NodeStatusBusy, Address: "3"}, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	status := pool.NodeStatusBusy
	count, err := storage.GetCountWithStatus(&status)
	assert.Nil(t, err)
	assert.Equal(t, count, 1)
}

// testStorage_Remove проверка удаления ноды
func testStorage_Remove(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	node := pool.Node{Status: pool.NodeStatusAvailable, Address: "1", CapabilitiesList: []capabilities.Capabilities{{"1": "2"}}}
	err := storage.Add(node, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	err = storage.Remove(node)
	assert.Nil(t, err)
	node, err = storage.GetByAddress(node.Address)
	assert.Error(t, err)
}

// testStorage_ReserveAvailable_Positive проверка резервирования ноды
func testStorage_ReserveAvailable_Positive(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	node := pool.Node{Status: pool.NodeStatusAvailable, Address: "1", CapabilitiesList: []capabilities.Capabilities{{"cap1": "val1"}}}
	err := storage.Add(node, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	expectedNode := pool.Node{Status: pool.NodeStatusAvailable, Address: "2", CapabilitiesList: []capabilities.Capabilities{{"cap1": "val1", "cap2": "val2"}}}
	err = storage.Add(expectedNode, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	node, err = storage.ReserveAvailable([]pool.Node{expectedNode})
	assert.Nil(t, err)
	assert.Equal(t, pool.NodeStatusReserved, node.Status, "Node not Reserved")
	assert.Equal(t, expectedNode.Address, node.Address, "Reserved unexpected node")
	node, err = storage.GetByAddress(node.Address)
	if err != nil {
		t.Fatal("Error get node, " + err.Error())
	}
	assert.Equal(t, pool.NodeStatusReserved, node.Status, "Node not Reserved")
}

// testStorage_ReserveAvailable_Negative проверка резервирования ноды, при условии отсутствия доступных нод
func testStorage_ReserveAvailable_Negative(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	node := pool.Node{Status: pool.NodeStatusBusy, Address: "qqqqqq", CapabilitiesList: []capabilities.Capabilities{{"1": "2"}}}
	err := storage.Add(node, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	node, err = storage.ReserveAvailable([]pool.Node{{Address: "qqqqqq"}})
	assert.Error(t, err)
}

// testStorage_SetAvailable проверка изменения статуса ноды на Available
func testStorage_SetAvailable(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	node := pool.Node{Status: pool.NodeStatusBusy, Address: "qqqqqq", CapabilitiesList: []capabilities.Capabilities{{"1": "2"}}}
	err := storage.Add(node, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	err = storage.SetAvailable(node)
	assert.Nil(t, err)
	node, err = storage.GetByAddress(node.Address)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	assert.Equal(t, pool.NodeStatusAvailable, node.Status, "Node not Available")
}

// testStorage_SetBusy проверка изменения статуса ноды на Busy
func testStorage_SetBusy(t *testing.T, p PrepareInterface) {
	t.Parallel()
	storage, deferFunc := p.CreateStorage()
	defer deferFunc()
	node := pool.Node{Status: pool.NodeStatusAvailable, Address: "qqqqqq", CapabilitiesList: []capabilities.Capabilities{{"1": "2"}}}
	err := storage.Add(node, 0)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	expectedSessionID := "newSessionId"
	err = storage.SetBusy(node, expectedSessionID)
	assert.Nil(t, err)
	node, err = storage.GetByAddress(node.Address)
	if err != nil {
		t.Fatal("Error add node, " + err.Error())
	}
	assert.Equal(t, pool.NodeStatusBusy, node.Status, "Node not Busy")
	assert.Equal(t, expectedSessionID, node.SessionID, "Not saved sessionID")
}
