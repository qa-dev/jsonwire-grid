package persistent

import (
	"fmt"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
)

type nodeHelperFactory struct{}

func (f *nodeHelperFactory) create(abstractClient jsonwire.ClientInterface) sessionsRemover {
	return &nodeHelper{client: abstractClient}
}

type nodeHelper struct {
	client jsonwire.ClientInterface
}

func (n *nodeHelper) removeAllSessions() (int, error) {
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
