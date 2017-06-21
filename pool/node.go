package pool

import "github.com/qa-dev/jsonwire-grid/pool/capabilities"

type NodeStatus string

const (
	NodeStatusAvailable NodeStatus = "available"
	NodeStatusReserved  NodeStatus = "reserved"
	NodeStatusBusy      NodeStatus = "busy"
)

type NodeType string

const (
	NodeTypePersistent NodeType = "persistent"
	NodeTypeKubernetes NodeType = "kubernetes"
)

type Node struct {
	Type             NodeType
	Address          string
	Status           NodeStatus
	SessionID        string
	Updated          int64
	Registered       int64
	CapabilitiesList []capabilities.Capabilities
}

func (n *Node) String() string {
	return "Node [" + n.Address + "]"
}

func NewNode(
	t NodeType,
	address string,
	status NodeStatus,
	sessionID string,
	updated int64,
	registered int64,
	capabilitiesList []capabilities.Capabilities,
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
