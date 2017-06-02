package jsonwire

import (
	"encoding/json"
	"net/http"
)

type ResponseStatus int

const (
	RESPONSE_STATUS_SUCCESS     ResponseStatus = 0
	RESPONSE_STATUS_UNKNOWN_ERR ResponseStatus = 13
)

type Response struct {
	SessionId *string         `json:"sessionId"`
	Status    ResponseStatus  `json:"status"`
	Value     json.RawMessage `json:"value"`
}

func NewResponse(sessionId *string, status ResponseStatus, value json.RawMessage) *Response {
	return &Response{sessionId, status, value}
}

func JsonResponse(w http.ResponseWriter, data Response, code int) (int, error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	body, err := json.Marshal(data)
	if err != nil {
		body = []byte(
			`{"sessionId": null, "status": 13}`)
	}

	return w.Write(body)
}
