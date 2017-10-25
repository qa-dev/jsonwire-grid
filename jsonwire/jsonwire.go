package jsonwire

import (
	"encoding/json"
)

type Message struct {
	SessionID string      `json:"sessionId"`
	Status    int         `json:"status"`
	Value     interface{} `json:"value"`
}

type NewSession struct {
	Message
	Value struct {
		SessionID string `json:"sessionId"`
	} `json:"value"`
}

type Sessions struct {
	Message
	Value []struct {
		ID           string          `json:"id"`
		Capabilities json.RawMessage `json:"capabilities"`
	} `json:"value"`
}

type Register struct {
	Class            string         `json:"class,omitempty"`
	Configuration    *Configuration `json:"configuration,omitempty"`
	CapabilitiesList []Capabilities `json:"capabilities,omitempty"` // selenium 3
}

type Capabilities map[string]interface{}

type Configuration struct {
	Id               string         `json:"id,omitempty"`
	Proxy            string         `json:"proxy,omitempty"`
	Role             string         `json:"role,omitempty"`
	Hub              string         `json:"hub,omitempty"`
	Port             int            `json:"port,omitempty"`
	RemoteHost       string         `json:"remoteHost,omitempty"`
	Host             string         `json:"host,omitempty"`
	MaxSession       int            `json:"maxSession,omitempty"`
	HubHost          string         `json:"hubHost,omitempty"`
	RegisterCycle    int            `json:"registerCycle,omitempty"`
	HubPort          int            `json:"hubPort,omitempty"`
	URL              string         `json:"url,omitempty"`
	Register         bool           `json:"register,omitempty"`
	CapabilitiesList []Capabilities `json:"capabilities,omitempty"` // selenium 2
}

type APIProxy struct {
	ID      string      `json:"id"`
	Request interface{} `json:"request"` //todo: пока не ясно зачем он нужен
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
}
