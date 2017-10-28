package capabilities

import "github.com/stretchr/testify/mock"

// ComparatorMock - mock of comparator.
type ComparatorMock struct {
	mock.Mock
}

// ComparatorMock - mock of compare capabilities.
func (c *ComparatorMock) Compare(desired Capabilities, available Capabilities) bool {
	args := c.Called(desired, available)
	return args.Bool(0)
}

// Register - mock of register capabilities for filtration.
func (c *ComparatorMock) Register(caps Capabilities) {
	_ = c.Called(caps)
}
