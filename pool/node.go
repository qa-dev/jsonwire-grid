package pool

type NodeStatus string

const (
	NodeStatusAvailable NodeStatus = "available"
	NodeStatusReserved  NodeStatus = "reserved"
	NodeStatusBusy      NodeStatus = "busy"
)

type NodeType string

const (
	NodeTypeRegular   NodeType = "regular"
	NodeTypeTemporary NodeType = "temporary"
)

type Node struct {
	Type             NodeType
	Address          string
	Status           NodeStatus
	SessionID        string
	Updated          int64
	Registered       int64
	CapabilitiesList []Capabilities
}

func (n *Node) String() string {
	return "Node [" + n.Address + "]"
}

type Capabilities map[string]interface{}

func NewNode(
	t NodeType,
	address string,
	status NodeStatus,
	sessionID string,
	updated int64,
	registered int64,
	capabilitiesList []Capabilities,
) *Node {
	return &Node{
		t,
		address,
		status,
		sessionID,
		updated,
		registered,
		capabilitiesList,
	}
}
