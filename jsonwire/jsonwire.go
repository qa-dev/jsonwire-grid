package jsonwire

import (
	"encoding/json"
)

type Message struct {
	SessionId string        `json:"sessionId"`
	Status    int           `json:"status"`
	Value     []interface{} `json:"value"`
}

type NewSession struct {
	Message
	Value struct {
		SessionId string `json:"sessionId"`
	} `json:"value"`
}

type Sessions struct {
	Message
	Value []struct {
		Id           string          `json:"id"`
		Capabilities json.RawMessage `json:"capabilities"`
	} `json:"value"`
}
