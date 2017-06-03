package jsonwire

import (
	"errors"
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
		err = errors.New(fmt.Sprintf("Can't get sessions, %s", err.Error()))
		return 0, err
	}
	if message.Status != 0 {
		err = errors.New(fmt.Sprintf("client.Sessions: Not succcess response status node, %s", n.client.Address()))
		return 0, err
	}
	// hasn't sessions
	countSessions := len(message.Value)
	if countSessions == 0 {
		return 0, nil
	}
	for _, session := range message.Value {
		message, err := n.client.CloseSession(session.Id)
		if err != nil {
			err = errors.New(fmt.Sprintf("Can't close session[%s], %s", session.Id, err.Error()))
			return 0, err
		}
		if message.Status != 0 {
			err = errors.New(fmt.Sprintf("client.CloseSession[%s]: Not succcess response node, %s", session.Id, n.client.Address()))
			return 0, err
		}
	}
	return countSessions, nil
}
