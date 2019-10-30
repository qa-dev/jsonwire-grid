package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
)

// RegisterNode - Receives requests to register nodes for persistent strategy.
type RegisterNode struct {
	Pool *pool.Pool
}

func (h *RegisterNode) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorMessage := "Can't read register request body: " + err.Error()
		log.Warning(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}
	err = r.Body.Close()
	if err != nil {
		log.Errorf("register/node: close request body, %v", err)
	}
	var register jsonwire.Register
	err = json.Unmarshal(body, &register)
	if err != nil {
		errorMessage := "Can't unmarshal register json:  " + err.Error()
		log.Warning(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}

	var capabilitiesList []jsonwire.Capabilities

	//todo: знаю что костыль, переделаю
	if len(register.Configuration.CapabilitiesList) > 0 {
		capabilitiesList = register.Configuration.CapabilitiesList
	} else {
		capabilitiesList = register.CapabilitiesList
	}
	poolCapabilitiesList := make([]capabilities.Capabilities, len(capabilitiesList))
	for i, value := range capabilitiesList {
		poolCapabilitiesList[i] = capabilities.Capabilities(value)
	}
	hostPort := register.Configuration.Host + ":" + strconv.Itoa(register.Configuration.Port)
	err = h.Pool.Add(
		hostPort,
		pool.NodeTypePersistent,
		hostPort,
		poolCapabilitiesList,
	)
	if err != nil {
		errorMessage := "Can't add node to pool:  " + err.Error()
		log.Warning(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write([]byte("ok"))
	if err != nil {
		log.Errorf("register/node: write response, %v", err)
	}
}
