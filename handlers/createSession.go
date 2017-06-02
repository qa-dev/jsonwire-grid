package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"jsonwire-grid/jsonwire"
	"jsonwire-grid/pool"
	"jsonwire-grid/proxy"
	"jsonwire-grid/selenium"
	"jsonwire-grid/wda"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type CreateSession struct {
	Pool *pool.Pool
}

func (h *CreateSession) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var capabilities map[string]Capabilities
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorMessage := "Error reading request: " + err.Error()
		log.Error(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}
	err = r.Body.Close()
	if err != nil {
		log.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rc := ioutil.NopCloser(bytes.NewReader(body))
	r.Body = rc
	log.Infof("requested session with params: %s", string(body))
	err = json.Unmarshal(body, &capabilities)
	if err != nil {
		errorMessage := "Error unmarshal json: " + err.Error()
		log.Error(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}
	desiredCapabilities, ok := capabilities["desiredCapabilities"]
	if !ok {
		errorMessage := "Not passed 'desiredCapabilities'"
		log.Error(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}
	poolCapabilities := pool.Capabilities(desiredCapabilities)
	tw, err := h.tryCreateSession(r, &poolCapabilities)
	if err != nil {
		http.Error(rw, "Can't create session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(tw.StatusCode)
	rw.Write(tw.Output)
}

func (h *CreateSession) tryCreateSession(r *http.Request, capabilities *pool.Capabilities) (*proxy.ResponseWriter, error) {
	//todo: если запрос отменить, все равно получение ноды будет повторяться, придумать как это предотвратить
	node, err := h.Pool.ReserveAvailableNode(*capabilities)
	if err != nil {
		return nil, errors.New("reserve node error: " + err.Error())
	}
	//todo: посылать в мониторинг событие, если вернулся не 0
	seleniumClient, err := createClient(node.Address, capabilities)
	if err != nil {
		return nil, errors.New("create Client error: " + err.Error())
	}
	seleniumNode := jsonwire.NewNode(seleniumClient)
	_, err = seleniumNode.RemoveAllSessions()
	if err != nil {
		log.Warn("Can't remove all sessions from node, go to next available node: " + node.String())
		h.Pool.Remove(node)
		return h.tryCreateSession(r, capabilities)
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
		log.Warn("Failure proxy request on node " + node.String() + ": " + string(tw.Output))
		h.Pool.Remove(node)
		return h.tryCreateSession(r, capabilities)
	}

	return tw, nil
}

func createClient(addr string, capabilities *pool.Capabilities) (jsonwire.ClientInterface, error) {
	if capabilities == nil {
		return nil, errors.New("capabilities must be not nil")
	}
	platformName := (*capabilities)["platformName"]
	switch platformName {
	case "WDA":
		return wda.NewClient(addr), nil
	default:
		return selenium.NewClient(addr), nil
	}
}
