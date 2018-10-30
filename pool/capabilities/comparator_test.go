package capabilities

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewComparator(t *testing.T) {
	t.Parallel()
	comp := NewComparator()
	assert.NotNil(t, comp)
}

func TestComparator_Compare(t *testing.T) {
	t.Parallel()
	var dataProvider = []struct {
		required       Capabilities
		available      Capabilities
		expectedResult bool
	}{
		{
			Capabilities{"browserName": "firefox"},
			Capabilities{"browserName": "firefox"},
			true,
		},
		{
			Capabilities{"browserName": "firefox"},
			Capabilities{"browserName": "chrome"},
			false,
		},
		{
			Capabilities{"browserName": "firefox"},
			Capabilities{"browserName": "firefox", "browserVersion": 1},
			true,
		},
		{
			Capabilities{"browserName": "firefox", "browserVersion": 1},
			Capabilities{"browserName": "firefox"},
			false,
		},
		{
			Capabilities{"browserName": "firefox", "myDogName": "petr"},
			Capabilities{"browserName": "firefox"},
			true,
		},
		{
			required:       Capabilities{"browserName": "firefox"},
			expectedResult: false,
		},
		{
			expectedResult: true,
		},
		{
			Capabilities{"platform": "ANY", "myDogName": "petr"},
			Capabilities{"platform": "LINUX-TORVALDS"},
			true,
		},
		{
			Capabilities{"platform": "ANY", "myDogName": "petr"},
			Capabilities{"not-defined-platform": "trololo"},
			true,
		},
	}

	comp := NewComparator()
	for i, test := range dataProvider {
		comp.Register(test.available)
		result := comp.Compare(test.required, test.available)
		assert.Equal(t, test.expectedResult, result, "Test #"+strconv.Itoa(i))
	}
}

func TestComparator_Register(t *testing.T) {
	t.Parallel()
	comp := NewComparator()
	comp.Register(Capabilities{"myNewSuperCapability": "anyValue", "myNewSuperCapability2": "anyValue"})
	assert.Contains(t, comp.registeredCaps, "myNewSuperCapability")
	assert.Contains(t, comp.registeredCaps, "myNewSuperCapability2")
}

func TestComparator_isRegistered_Positive(t *testing.T) {
	t.Parallel()
	comp := NewComparator()
	comp.registeredCaps["myNewSuperCapabiluty"] = struct{}{}
	assert.True(t, comp.isRegistered("myNewSuperCapabiluty"))
}

func TestComparator_isRegistered_Negative(t *testing.T) {
	t.Parallel()
	comp := NewComparator()
	assert.False(t, comp.isRegistered("myNewSuperCapabiluty"))
}
