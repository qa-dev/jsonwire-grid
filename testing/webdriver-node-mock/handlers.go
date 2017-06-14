package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"math/rand"
	"net/http"
	"os/exec"
	"regexp"
	"time"
)

// getSessions return list active sessions
func getSessions(rw http.ResponseWriter, r *http.Request) {
	sessions := &jsonwire.Sessions{}
	if currentSessionID != "" {
		sessions.Value = []struct {
			Id           string          `json:"id"`
			Capabilities json.RawMessage `json:"capabilities"`
		}{
			{Id: currentSessionID, Capabilities: nil},
		}
	}

	err := json.NewEncoder(rw).Encode(sessions)
	if err != nil {
		err = errors.New("Get sessions error, " + err.Error())
		log.Error(err)
		json.NewEncoder(rw).Encode(&jsonwire.Message{Value: err.Error(), Status: int(jsonwire.RESPONSE_STATUS_UNKNOWN_ERR)})
	}
}

// createSession create new session
func createSession(rw http.ResponseWriter, r *http.Request) {
	if maxDuration > 0 {
		rand.Seed(time.Now().UTC().UnixNano())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(maxDuration)))
	}
	rw.Header().Set("Accept", "application/json")
	rw.Header().Set("Accept-charset", "utf-8")

	if r.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if currentSessionID != "" {
		errorMassage := "Session already exists"
		log.Error(errorMassage)
		rw.WriteHeader(http.StatusInternalServerError)
		responseMessage := &jsonwire.Message{}
		responseMessage.Status = int(jsonwire.RESPONSE_STATUS_UNKNOWN_ERR)
		responseMessage.Value = errorMassage
		json.NewEncoder(rw).Encode(responseMessage)
		return
	}

	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatalf("Can't generate uuid, %s", err)
	}
	currentSessionID = string(out[:len(out)-1]) // cut end of line char
	json.NewEncoder(rw).Encode(&jsonwire.Message{SessionId: currentSessionID})
}

// useSession mocks any action as open page, click, close session
func useSession(rw http.ResponseWriter, r *http.Request) {
	if maxDuration > 0 {
		rand.Seed(time.Now().UTC().UnixNano())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(maxDuration)))
	}
	rw.Header().Set("Accept", "application/json")
	rw.Header().Set("Accept-charset", "utf-8")

	re := regexp.MustCompile(".*/session/([^/]+)(?:/([^/]+))?")
	parsedUrl := re.FindStringSubmatch(r.URL.Path)
	if len(parsedUrl) != 3 {
		errorMessage := "url [" + r.URL.Path + "] parsing error"
		log.Infof(errorMessage)
		http.Error(rw, errorMessage, http.StatusBadRequest)
		return
	}
	sessionId := re.FindStringSubmatch(r.URL.Path)[1]
	responseMessage := &jsonwire.Message{SessionId: sessionId}
	if sessionId != currentSessionID {
		errorMassage := fmt.Sprintf("sessionID '%s' not found", sessionId)
		log.Error(errorMassage)
		rw.WriteHeader(http.StatusNotFound)
		responseMessage.Status = int(jsonwire.RESPONSE_STATUS_UNKNOWN_ERR)
		responseMessage.Value = errorMassage
		json.NewEncoder(rw).Encode(responseMessage)
		return
	}
	if parsedUrl[2] == "" && r.Method == http.MethodDelete {
		currentSessionID = ""
	}
	json.NewEncoder(rw).Encode(responseMessage)
}
