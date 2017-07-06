package jsonwire

import (
	"github.com/stretchr/testify/mock"
)

type ClientFactoryMock struct {
	mock.Mock
}

func (cf *ClientFactoryMock) Create(address string) ClientInterface {
	args := cf.Called(address)
	return args.Get(0).(ClientInterface)
}

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) Status() (*Message, error) {
	args := c.Called()
	return args.Get(0).(*Message), args.Error(1)
}

func (c *ClientMock) Sessions() (*Sessions, error) {
	args := c.Called()
	return args.Get(0).(*Sessions), args.Error(1)
}

func (c *ClientMock) CloseSession(sessionID string) (*Message, error) {
	args := c.Called(sessionID)
	return args.Get(0).(*Message), args.Error(1)
}

func (c *ClientMock) Address() string {
	args := c.Called()
	return args.String(0)
}
