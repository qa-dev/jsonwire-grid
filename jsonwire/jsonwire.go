package jsonwire

import (
	"encoding/json"
)

// Message - common protocol message structure.
type Message struct {
	SessionID string      `json:"sessionId"`
	Status    int         `json:"status"`
	Value     interface{} `json:"value"`
}

// NewSession - message structure for the creation of a new session
type NewSession struct {
	Message
	Value struct {
		SessionID string `json:"sessionId"`
	} `json:"value"`
}

// Sessions - message structure for sessions list.
type Sessions struct {
	Message
	Value []struct {
		ID           string          `json:"id"`
		Capabilities json.RawMessage `json:"capabilities"`
	} `json:"value"`
}

// Register - message structure for registration new node.
type Register struct {
	Class            json.RawMessage `json:"class,omitempty"`
	Configuration    *Configuration  `json:"configuration,omitempty"`
	CapabilitiesList []Capabilities  `json:"capabilities,omitempty"` // selenium 3
	Description		 string			 `json:"description,omitempty"`
	Name		 	 string			 `json:"name,omitempty"`
}

// Capabilities - structure of supported capabilities.
type Capabilities map[string]interface{}

// Configuration - structure of node configuration.
type Configuration struct {
	ID         string `json:"id,omitempty"`
	Proxy      string `json:"proxy,omitempty"`
	Role       string `json:"role,omitempty"`
	Hub        string `json:"hub,omitempty"`
	Port       int    `json:"port,omitempty"`
	RemoteHost string `json:"remoteHost,omitempty"`
	Host       string `json:"host,omitempty"`
	MaxSession int    `json:"maxSession,omitempty"`
	HubHost    string `json:"hubHost,omitempty"`
	RegisterCycle    int            `json:"registerCycle,omitempty"`
	HubPort          int            `json:"hubPort,omitempty"`
	URL              string         `json:"url,omitempty"`
	Register         bool           `json:"register,omitempty"`
	CapabilitiesList []Capabilities `json:"capabilities,omitempty"` // selenium 2
}

// APIProxy - message structure for node check.
type APIProxy struct {
	ID      string      `json:"id"`
	Request interface{} `json:"request"` //todo: пока не ясно зачем он нужен
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
}
