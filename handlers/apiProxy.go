package handlers

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid/pool"
)

type ApiProxy struct {
	Pool *pool.Pool
}

func (h *ApiProxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-type", "application/json")

	id := r.URL.Query().Get("id")
	nodeUrl, err := url.Parse(id) //todo: обработка ошибок

	if err != nil {
		errorMessage := "Error get 'id' from url: " + r.URL.String()
		log.Warning(errorMessage)
		rw.Write([]byte(`{"id": "", "request": {}, "msg": "` + errorMessage + `": false}`))
		return
	}

	node, err := h.Pool.GetNodeByAddress(nodeUrl.Host)

	//todo: хардкод для ткста, сделать нормальный респонз обжекты
	if node == nil || err != nil {
		errorMessage := "api/proxy: Can't get node"
		if err != nil {
			errorMessage = err.Error()
			log.Warning(errorMessage)
		}
		rw.Write([]byte(`{"msg": "Cannot find proxy with ID =` + id + ` in the registry: ` + errorMessage + `", "success": false}`))
	} else {
		rw.Write([]byte(`{"id": "", "request": {}, "msg": "proxy found !", "success": true}`))
	}
}
