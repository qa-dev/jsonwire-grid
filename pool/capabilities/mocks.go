package capabilities

import "github.com/stretchr/testify/mock"

type ComparatorMock struct {
	mock.Mock
}

func (c *ComparatorMock) Compare(desired Capabilities, available Capabilities) bool {
	args := c.Called(desired, available)
	return args.Bool(0)
}

func (c *ComparatorMock) Register(caps Capabilities) {
	_ = c.Called(caps)
}
