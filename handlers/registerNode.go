package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"jsonwire-grid/pool"
)

type RegisterNode struct {
	Pool *pool.Pool
}

type registerJson struct {
	Class            string             `json:"class"`
	Configuration    *configurationJson `json:"configuration"`
	CapabilitiesList []Capabilities     `json:"capabilities"` // selenium 3
}

type Capabilities map[string]interface{}

type configurationJson struct {
	Proxy            string
	Role             string
	Hub              string
	Port             int
	RemoteHost       string
	Host             string
	MaxSession       int
	HubHost          string
	RegisterCycle    int
	HubPort          int
	Url              string
	Register         bool
	CapabilitiesList []Capabilities `json:"capabilities"` // selenium 2
}

func (h *RegisterNode) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorMessage := "Can't read register request body: " + err.Error()
		log.Warning(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	var register registerJson
	err = json.Unmarshal(body, &register)
	if err != nil {
		errorMessage := "Can't unmarshal register json:  " + err.Error()
		log.Warning(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}

	var capabilitiesList []Capabilities

	//todo: знаю что костыль, переделаю
	if len(register.Configuration.CapabilitiesList) > 0 {
		capabilitiesList = register.Configuration.CapabilitiesList
	} else {
		capabilitiesList = register.CapabilitiesList
	}
	poolCapabilitiesList := make([]pool.Capabilities, len(capabilitiesList))
	for i, value := range capabilitiesList {
		poolCapabilitiesList[i] = pool.Capabilities(value)
	}
	err = h.Pool.Add(
		pool.NodeTypeRegular,
		register.Configuration.Host+":"+strconv.Itoa(register.Configuration.Port),
		poolCapabilitiesList,
	)
	if err != nil {
		errorMessage := "Can't add node to pool:  " + err.Error()
		log.Warning(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("ok"))
}
