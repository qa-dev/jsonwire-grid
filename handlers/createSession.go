package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/qa-dev/jsonwire-grid/proxy"
)

// CreateSession - Receives requests to create session.
type CreateSession struct {
	Pool          *pool.Pool
	ClientFactory jsonwire.ClientFactoryInterface
}

func (h *CreateSession) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var caps map[string]jsonwire.Capabilities
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorMessage := "Error reading request: " + err.Error()
		log.Error(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}
	err = r.Body.Close()
	if err != nil {
		log.Errorf("create/session: close request body, %v", err)
	}
	rc := ioutil.NopCloser(bytes.NewReader(body))
	r.Body = rc
	log.Infof("requested session with params: %s", string(body))
	err = json.Unmarshal(body, &caps)
	if err != nil {
		errorMessage := "Error unmarshal json: " + err.Error()
		log.Error(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}
	desiredCapabilities, ok := caps["desiredCapabilities"]
	if !ok {
		errorMessage := "Not passed 'desiredCapabilities'"
		log.Error(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}
	poolCapabilities := capabilities.Capabilities(desiredCapabilities)
	tw, err := h.tryCreateSession(r, &poolCapabilities)
	if err != nil {
		http.Error(rw, "Can't create session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(tw.StatusCode)
	_, err = rw.Write(tw.Output)
	if err != nil {
		log.Errorf("create/session: write response, %v", err)
	}
}

func (h *CreateSession) tryCreateSession(r *http.Request, capabilities *capabilities.Capabilities) (*proxy.ResponseWriter, error) {
	select {
	case <-r.Context().Done():
		err := errors.New("Request canceled by client, " + r.Context().Err().Error())
		return nil, err
	default:
	}

	node, err := h.Pool.ReserveAvailableNode(*capabilities)
	if err != nil {
		return nil, errors.New("reserve node error: " + err.Error())
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   node.Address,
	})
	transport := proxy.NewCreateSessionTransport(h.Pool, node)
	reverseProxy.Transport = transport
	tw := proxy.NewResponseWriter()
	reverseProxy.ServeHTTP(tw, r)

	if !transport.IsSuccess {
		err = h.Pool.CleanUpNode(node)
		if err != nil {
			log.Errorf("fail cleanUp node on create session failure, %s", err)
		}
		return nil, fmt.Errorf("failure proxy request on node %s: %s msg: %s", node, tw.Output, transport.Error)
	}

	return tw, nil
}
