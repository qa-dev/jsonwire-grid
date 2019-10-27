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
	// A unique key, by which we understand how to find this object in the outer world + for not adding the second time the same thing.
	// The value may depend on the strategy:
	// - for constant nodes ip: port
	// - for temporary pod.name
	Key              string                      `json:"key"`
	Type             NodeType                    `json:"type"`
	Address          string                      `json:"address"`
	Status           NodeStatus                  `json:"status"`
	SessionID        string                      `json:"session_id"`
	Updated          int64                       `json:"updated"`
	Registered       int64                       `json:"registered"`
	CapabilitiesList []capabilities.Capabilities `json:"capabilities_list"`
}

func (n *Node) String() string {
	return "Node [" + n.Key + "]"
}

func NewNode(
	key string,
	t NodeType,
	address string,
	status NodeStatus,
	sessionID string,
	updated int64,
	registered int64,
	capabilitiesList []capabilities.Capabilities,
) *Node {
	return &Node{
		key,
		t,
		address,
		status,
		sessionID,
		updated,
		registered,
		capabilitiesList,
	}
}
