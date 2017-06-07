package jsonwire

import (
	"encoding/json"
)

type Message struct {
	SessionId string        `json:"sessionId"`
	Status    int           `json:"status"`
	Value     interface{} `json:"value"`
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

type Register struct {
	Class            string             `json:"class"`
	Configuration    *Configuration `json:"configuration"`
	CapabilitiesList []Capabilities     `json:"capabilities"` // selenium 3
}

type Capabilities map[string]interface{}

type Configuration struct {
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

type ApiProxy struct {
	ID string `json:"id"`
	Request interface{} `json:"request"` //todo: пока не ясно зачем он нужен
	Msg string `json:"msg"`
	Success bool `json:"success"`
}