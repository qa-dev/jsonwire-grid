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
		return nil, errors.New("round trip to node: " + err.Error())
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("read node response: " + err.Error())
	}
	err = response.Body.Close()
	if err != nil {
		return nil, errors.New("close response body: " + err.Error())
	}

	var message jsonwire.NewSession
	err = json.Unmarshal(b, &message)
	if err != nil {
		return nil, errors.New("read body with sessionID: " + err.Error())
	}
	var sessionID string
	switch {
	case message.SessionID != "":
		sessionID = message.SessionID
	case message.Value.SessionID != "":
		sessionID = message.Value.SessionID
	default:
		return nil, fmt.Errorf("session not created, response: %s", string(b))
	}
	log.Infof("register SessionID: %s on node %s", sessionID, t.node.Address)
	err = t.pool.RegisterSession(t.node, sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessionId not registred in storage: %s", sessionID)
	}
	response.Body = ioutil.NopCloser(bytes.NewReader(b))
	t.IsSuccess = true
	return response, err
}
