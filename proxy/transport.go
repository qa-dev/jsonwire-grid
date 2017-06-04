package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"errors"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"github.com/qa-dev/jsonwire-grid/pool"
)

type CreateSessionTransport struct {
	pool      *pool.Pool
	node      *pool.Node
	IsSuccess bool
}

func NewCreateSessionTransport(pool *pool.Pool, node *pool.Node) *CreateSessionTransport {
	return &CreateSessionTransport{pool: pool, node: node}
}

func (t *CreateSessionTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		err = errors.New("Can't round trip to node: " + err.Error())
		return nil, err
	}

	// small quantity ebamistic-magic
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New("Can't read node response: " + err.Error())
		return nil, err
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}
	body := ioutil.NopCloser(bytes.NewReader(b))
	response.Body = body

	var message jsonwire.NewSession
	err = json.Unmarshal(b, &message)
	if err != nil {
		response.Body.Close()
		err = errors.New("Cant read sessionID: " + err.Error())
		return nil, err
	}
	var sessionID string
	switch {
	case message.SessionId != "":
		sessionID = message.SessionId
	case message.Value.SessionId != "":
		sessionID = message.Value.SessionId
	default:
		response.Body.Close()
		return nil, errors.New(fmt.Sprintf("Session not created, response: %s", string(b)))
	}
	log.Infof("register SessionID: %s on node %s", sessionID, t.node.Address)
	err = t.pool.RegisterSession(t.node, sessionID)
	if err != nil {
		response.Body.Close()
		log.Errorf("sessionId not registred in storage: %s", sessionID)
		return nil, err
	}
	t.IsSuccess = true
	return response, err
}
