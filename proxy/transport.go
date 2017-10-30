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
	Error     error
}

func NewCreateSessionTransport(pool *pool.Pool, node *pool.Node) *CreateSessionTransport {
	return &CreateSessionTransport{pool: pool, node: node}
}

func (t *CreateSessionTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var err error
	defer func() { t.Error = err }() // dirty hack, for get error from round trip
	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		err = errors.New("round trip to node: " + err.Error())
		return nil, err
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New("read node response: " + err.Error())
		return nil, err
	}
	err = response.Body.Close()
	if err != nil {
		err = errors.New("close response body: " + err.Error())
		return nil, err
	}

	var message jsonwire.NewSession
	err = json.Unmarshal(b, &message)
	if err != nil {
		err = errors.New("read body with sessionID: " + err.Error())
		return nil, err
	}
	var sessionID string
	switch {
	case message.SessionID != "":
		sessionID = message.SessionID
	case message.Value.SessionID != "":
		sessionID = message.Value.SessionID
	default:
		err = fmt.Errorf("session not created, response: %s", string(b))
		return nil, err
	}
	log.Infof("register SessionID: %s on node %s", sessionID, t.node.Address)
	err = t.pool.RegisterSession(t.node, sessionID)
	if err != nil {
		err = fmt.Errorf("sessionId not registred in storage: %s", sessionID)
		return nil, err
	}
	response.Body = ioutil.NopCloser(bytes.NewReader(b))
	t.IsSuccess = true
	return response, nil
}
