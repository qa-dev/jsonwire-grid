package persistent

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"errors"
)

func TestNodeHelperFactory_create(t *testing.T) {
	nhf := new(nodeHelperFactory)
	nodeHelper := nhf.create(nil)
	assert.NotNil(t, nodeHelper)
}

func TestNodeHelper_removeAllSessions_Positive_NothingToRemove(t *testing.T) {
	cm := new(jsonwire.ClientMock)
	nodeHelper := &nodeHelper{cm}
	message := new(jsonwire.Sessions)
	cm.On("Sessions").Return(message, nil)
	_, err := nodeHelper.removeAllSessions()
	assert.Nil(t, err)
}

func TestNodeHelper_removeAllSessions_Positive_SuccessRemove(t *testing.T) {
	cm := new(jsonwire.ClientMock)
	nodeHelper := &nodeHelper{cm}
	sessions := new(jsonwire.Sessions)
	sessions.Value = []struct {
		ID           string          `json:"id"`
		Capabilities json.RawMessage `json:"capabilities"`
	}{
		{ID: "lrololo", Capabilities: nil},
	}
	cm.On("Sessions").Return(sessions, nil)
	message := new(jsonwire.Message)
	cm.On("CloseSession", mock.AnythingOfType("string")).Return(message, nil)
	_, err := nodeHelper.removeAllSessions()
	assert.Nil(t, err)
}

func TestNodeHelper_removeAllSessions_Negative_Sessions_Error(t *testing.T) {
	cm := new(jsonwire.ClientMock)
	nodeHelper := &nodeHelper{cm}
	cm.On("Sessions").Return(new(jsonwire.Sessions), errors.New("Err"))
	_, err := nodeHelper.removeAllSessions()
	assert.NotNil(t, err)
}

func TestNodeHelper_removeAllSessions_Negative_Sessions_MessageStatusNotOk(t *testing.T) {
	cm := new(jsonwire.ClientMock)
	nodeHelper := &nodeHelper{cm}
	sessions := new(jsonwire.Sessions)
	sessions.Status = 99999
	cm.On("Sessions").Return(new(jsonwire.Sessions), errors.New("Err"))
	_, err := nodeHelper.removeAllSessions()
	assert.NotNil(t, err)
}

func TestNodeHelper_removeAllSessions_Negative_CloseSession_Error(t *testing.T) {
	cm := new(jsonwire.ClientMock)
	nodeHelper := &nodeHelper{cm}
	sessions := new(jsonwire.Sessions)
	sessions.Value = []struct {
		ID           string          `json:"id"`
		Capabilities json.RawMessage `json:"capabilities"`
	}{
		{ID: "lrololo", Capabilities: nil},
	}
	cm.On("Sessions").Return(sessions, nil)
	message := new(jsonwire.Message)
	cm.On("CloseSession", mock.AnythingOfType("string")).Return(message, errors.New("Err"))
	_, err := nodeHelper.removeAllSessions()
	assert.NotNil(t, err)
}

func TestNodeHelper_removeAllSessions_Negative_CloseSession_MessageStatusNotOk(t *testing.T) {
	cm := new(jsonwire.ClientMock)
	nodeHelper := &nodeHelper{cm}
	sessions := new(jsonwire.Sessions)
	sessions.Value = []struct {
		ID           string          `json:"id"`
		Capabilities json.RawMessage `json:"capabilities"`
	}{
		{ID: "lrololo", Capabilities: nil},
	}
	cm.On("Sessions").Return(sessions, nil)
	message := new(jsonwire.Message)
	message.Status = 999999
	cm.On("CloseSession", mock.AnythingOfType("string")).Return(message, errors.New("Err"))
	_, err := nodeHelper.removeAllSessions()
	assert.NotNil(t, err)
}

