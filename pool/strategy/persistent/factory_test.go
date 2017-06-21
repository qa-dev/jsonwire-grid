package persistent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStrategyFactory_Create(t *testing.T) {
	f := StrategyFactory{}
	s, err := f.Create(nil, nil, nil)
	assert.NotNil(t, s)
	assert.Nil(t, err)
}
