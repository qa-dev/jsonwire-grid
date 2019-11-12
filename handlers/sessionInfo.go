package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/qa-dev/jsonwire-grid/pool"
)

// SessionInfo - Returns a session info (node address, status, etc)
type SessionInfo struct {
	Pool *pool.Pool
}

type sessionInfoResponse struct {
	NodeAddress string          `json:"node_address"`
	NodeType    pool.NodeType   `json:"node_type"`
	Status      pool.NodeStatus `json:"node_status"`
	SessionID   string          `json:"session_id"`
	Updated     int64           `json:"updated"`
	Registered  int64           `json:"registered"`
}

func (h *SessionInfo) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("sessionid")
	if sessionId == "" {
		http.Error(rw, fmt.Sprint("session id must be specified,"), http.StatusBadRequest)
		return
	}
	node, err := h.Pool.GetNodeBySessionID(sessionId)
	if err != nil {
		http.Error(rw, fmt.Sprint("trying to get a session data,", err), http.StatusInternalServerError)
		return
	}

	resp := sessionInfoResponse{
		node.Address,
		node.Type,
		node.Status,
		node.SessionID,
		node.Updated,
		node.Registered,
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, fmt.Sprint("trying to encode a response,", err), http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(respJSON)
	if err != nil {
		log.Error("session/info: write a response,", err)
	}
}
