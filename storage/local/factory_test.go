package local

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qa-dev/jsonwire-grid/config"
)

func TestFactory_Create_Positive(t *testing.T) {
	f := Factory{}
	storage, err := f.Create(config.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}
