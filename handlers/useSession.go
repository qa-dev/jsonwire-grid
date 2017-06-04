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
	Pool *pool.Pool
}

func (h *UseSession) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(".*/session/([^/]+)(?:/([^/]+))?")
	parsedUrl := re.FindStringSubmatch(r.URL.Path)
	if len(parsedUrl) != 3 {
		errorMessage := "url [" + r.URL.Path + "] parsing error"
		log.Infof(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}
	sessionId := re.FindStringSubmatch(r.URL.Path)[1]
	targetNode, err := h.Pool.GetNodeBySessionId(sessionId)
	if targetNode == nil || err != nil {
		errorMessage := ""
		if err != nil {
			errorMessage = err.Error()
		}
		// посылаем сообщение о том что сессия не найдена
		errorMessage = "session " + sessionId + " not found in node pool: " + errorMessage
		log.Infof(errorMessage)
		http.Error(rw, errorMessage, http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   targetNode.Address,
	})
	proxy.ServeHTTP(rw, r)

	//// todo: заговнокодим пока не появилось понимание как лучше сделать ибо сраный rest
	if parsedUrl[2] == "" && r.Method == http.MethodDelete {
		err := h.Pool.CleanUpNode(targetNode)
		if err != nil {
			log.Error("Clanup node status error: " + err.Error())
		}
	}
}
