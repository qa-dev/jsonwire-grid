package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/qa-dev/jsonwire-grid/jsonwire"
)

var constResponse = RandStringRunes(10000)

// status return current status
func status(rw http.ResponseWriter, r *http.Request) {
	sessions := &jsonwire.Message{}
	err := json.NewEncoder(rw).Encode(sessions)
	if err != nil {
		err = errors.New("Get sessions error, " + err.Error())
		log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// getSessions return list active sessions
func getSessions(rw http.ResponseWriter, r *http.Request) {
	sessions := &jsonwire.Sessions{}
	if currentSessionID != "" {
		sessions.Value = []struct {
			ID           string          `json:"id"`
			Capabilities json.RawMessage `json:"capabilities"`
		}{
			{ID: currentSessionID, Capabilities: nil},
		}
	}

	err := json.NewEncoder(rw).Encode(sessions)
	if err != nil {
		err = errors.New("Get sessions error, " + err.Error())
		log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// createSession create new session
func createSession(rw http.ResponseWriter, r *http.Request) {
	if maxDuration > 0 {
		rand.Seed(time.Now().UTC().UnixNano())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(maxDuration)))
	}
	rw.Header().Set("Accept", "application/json")
	rw.Header().Set("Content-type", "application/json; charset=utf-8")
	rw.Header().Set("Accept-charset", "utf-8")

	if r.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	responseMessage := new(jsonwire.Message)
	switch {
	case currentSessionID != "": // trying create session on busy node
		errorMassage := "Session already exists"
		log.Error(errorMassage)
		rw.WriteHeader(http.StatusInternalServerError)
		responseMessage.Status = int(jsonwire.ResponseStatusUnknownErr)
		responseMessage.Value = errorMassage
	default:
		currentSessionID = uuid.NewV4().String()
		responseMessage.SessionID = currentSessionID
	}

	err := json.NewEncoder(rw).Encode(responseMessage)
	if err != nil {
		http.Error(rw, "Error encode response to json, rawMessage: "+fmt.Sprintf("%+v", responseMessage), http.StatusInternalServerError)
	}
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
	responseMessage := &jsonwire.Message{SessionID: sessionId}

	switch {
	case sessionId != currentSessionID: // client requested unknown session id
		errorMassage := fmt.Sprintf("sessionID '%s' not found", sessionId)
		log.Error(errorMassage)
		rw.WriteHeader(http.StatusNotFound)
		responseMessage.Status = int(jsonwire.ResponseStatusUnknownErr)
		responseMessage.Value = errorMassage
	case parsedUrl[2] == "" && r.Method == http.MethodDelete: // session closed by client
		currentSessionID = ""
	default:
		responseMessage.Value = constResponse
	}
	err := json.NewEncoder(rw).Encode(responseMessage)
	if err != nil {
		http.Error(rw, "Error encode response to json, rawMessage: "+fmt.Sprintf("%+v", responseMessage), http.StatusInternalServerError)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
