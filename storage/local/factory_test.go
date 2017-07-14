package local

import (
	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFactory_Create_Positive(t *testing.T) {
	f := Factory{}
	storage, err := f.Create(config.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}
