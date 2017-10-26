package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid/pool"
)

type UseSession struct {
	Pool  *pool.Pool
	Cache *pool.Cache
}

func (h *UseSession) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(".*/session/([^/]+)(?:/([^/]+))?")
	parsedURL := re.FindStringSubmatch(r.URL.Path)
	if len(parsedURL) != 3 {
		errorMessage := "url [" + r.URL.Path + "] parsing error"
		log.Infof(errorMessage)
		http.Error(rw, errorMessage, http.StatusBadRequest)
		return
	}
	sessionID := re.FindStringSubmatch(r.URL.Path)[1]
	targetNode, ok := h.Cache.Get(sessionID)
	var err error
	if !ok {
		targetNode, err = h.Pool.GetNodeBySessionID(sessionID)
		if err != nil {
			errorMessage := "session " + sessionID + " not found in node pool: " + err.Error()
			log.Infof(errorMessage)
			http.Error(rw, errorMessage, http.StatusNotFound)
			return
		}
		h.Cache.Set(sessionID, targetNode)
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   targetNode.Address,
	})
	proxy.ServeHTTP(rw, r)

	//// todo: заговнокодим пока не появилось понимание как лучше сделать ибо сраный rest
	if parsedURL[2] == "" && r.Method == http.MethodDelete {
		err := h.Pool.CleanUpNode(targetNode)
		if err != nil {
			log.Error("Clanup node status error: " + err.Error())
		}
	}
}
