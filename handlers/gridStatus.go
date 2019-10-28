package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/qa-dev/jsonwire-grid/pool"
)

// GridStatus - Returns a status.
type GridStatus struct {
	Pool   *pool.Pool
	Config config.Config
}

type response struct {
	NodeList []pool.Node   `json:"node_list"`
	Config   config.Config `json:"config"`
}

func (h *GridStatus) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	nodeList, err := h.Pool.GetAll()
	if err != nil {
		http.Error(rw, fmt.Sprint("trying to get a node list from pool,", err), http.StatusInternalServerError)
		return
	}

	resp := response{NodeList: nodeList, Config: h.Config}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, fmt.Sprint("trying to encode a response,", err), http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(respJSON)
	if err != nil {
		log.Error("grid/status: write a response,", err)
	}
}
