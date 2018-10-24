package handlers

import (
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid/pool"
)

// APIProxy - responds with information about whether the node is registered.
type APIProxy struct {
	Pool *pool.Pool
}

func (h *APIProxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-type", "application/json")

	id := r.URL.Query().Get("id")
	nodeURL, err := url.Parse(id)

	if err != nil {
		errorMessage := "Error get 'id' from url: " + r.URL.String()
		log.Warning(errorMessage)
		_, err = rw.Write([]byte(`{"id": "", "request": {}, "msg": "` + errorMessage + `": false}`))
		if err != nil {
			log.Errorf("api/proxy: write response, %v", err)
		}
		return
	}

	_, err = h.Pool.GetNodeByAddress(nodeURL.Host)

	//todo: хардкод для ткста, сделать нормальный респонз обжекты
	if err != nil {
		errorMessage := "api/proxy: Can't get node, " + err.Error()
		log.Warning(errorMessage)
		_, err = rw.Write([]byte(`{"msg": "Cannot find proxy with ID =` + id + ` in the registry: ` + errorMessage + `", "success": false}`))
	} else {
		_, err = rw.Write([]byte(`{"id": "", "request": {}, "msg": "proxy found !", "success": true}`))
	}

	if err != nil {
		log.Errorf("api/proxy: write response, %v", err)
	}
}
