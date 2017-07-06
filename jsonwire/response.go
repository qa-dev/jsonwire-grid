package jsonwire

import (
	"encoding/json"
	"net/http"
)

type ResponseStatus int

const (
	ResponseStatusSuccess    ResponseStatus = 0
	ResponseStatusUnknownErr ResponseStatus = 13
)

type Response struct {
	SessionID *string         `json:"sessionId"`
	Status    ResponseStatus  `json:"status"`
	Value     json.RawMessage `json:"value"`
}

func NewResponse(sessionID *string, status ResponseStatus, value json.RawMessage) *Response {
	return &Response{sessionID, status, value}
}

func JSONResponse(w http.ResponseWriter, data Response, code int) (int, error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	body, err := json.Marshal(data)
	if err != nil {
		body = []byte(
			`{"sessionId": null, "status": 13}`)
	}

	return w.Write(body)
}
