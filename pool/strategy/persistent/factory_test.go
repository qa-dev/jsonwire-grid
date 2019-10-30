package persistent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrategyFactory_Create(t *testing.T) {
	f := StrategyFactory{}
	s, err := f.Create(nil, nil, nil)
	assert.NotNil(t, s)
	assert.Nil(t, err)
}
