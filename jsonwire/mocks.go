package jsonwire

import (
	"github.com/stretchr/testify/mock"
)

// ClientFactoryMock - mock of factory of client.
type ClientFactoryMock struct {
	mock.Mock
}

// Create - mock of create session.
func (cf *ClientFactoryMock) Create(address string) ClientInterface {
	args := cf.Called(address)
	return args.Get(0).(ClientInterface)
}

// ClientMock - mock of client.
type ClientMock struct {
	mock.Mock
}

// Health - mock of healthcheck.
func (c *ClientMock) Health() (*Message, error) {
	args := c.Called()
	return args.Get(0).(*Message), args.Error(1)
}

// Sessions - mock of sessions list.
func (c *ClientMock) Sessions() (*Sessions, error) {
	args := c.Called()
	return args.Get(0).(*Sessions), args.Error(1)
}

// CloseSession - mock of close session.
func (c *ClientMock) CloseSession(sessionID string) (*Message, error) {
	args := c.Called(sessionID)
	return args.Get(0).(*Message), args.Error(1)
}

// Address - mock of get node address.
func (c *ClientMock) Address() string {
	args := c.Called()
	return args.String(0)
}
