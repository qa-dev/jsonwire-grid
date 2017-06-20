package persistent

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/qa-dev/jsonwire-grid/config"
)

func TestStrategyFactory_Create(t *testing.T) {
	f := StrategyFactory{}
	s, err := f.Create(config.Config{}, new(StorageMock))
	assert.NotNil(t, s)
	assert.Nil(t, err)
}
