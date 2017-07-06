package jsonwire

import (
	"fmt"
)

type Node struct {
	client ClientInterface
}

func NewNode(abstractClient ClientInterface) *Node {
	return &Node{client: abstractClient}
}

func (n *Node) RemoveAllSessions() (int, error) {
	message, err := n.client.Sessions()
	if err != nil {
		return 0, fmt.Errorf("Can't get sessions, %s", err.Error())
	}
	if message.Status != 0 {
		return 0, fmt.Errorf("client.Sessions: Not succcess response status node, %s", n.client.Address())
	}
	// hasn't sessions
	countSessions := len(message.Value)
	if countSessions == 0 {
		return 0, nil
	}
	for _, session := range message.Value {
		message, err := n.client.CloseSession(session.ID)
		if err != nil {
			return 0, fmt.Errorf("Can't close session[%s], %s", session.ID, err.Error())
		}
		if message.Status != 0 {
			return 0, fmt.Errorf("client.CloseSession[%s]: Not succcess response node, %s", session.ID, n.client.Address())
		}
	}
	return countSessions, nil
}
