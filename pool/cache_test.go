package pool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	c := NewCache(time.Second)
	assert.NotNil(t, c)
}

func TestCache_Set_ReturnsNode(t *testing.T) {
	c := NewCache(time.Second)
	key := "1"
	node := new(Node)
	c.Set(key, node)
	assert.Equal(t, node, c.storage[key].node)
}

func TestCache_Get_NodeExists_ReturnsNodeTrue(t *testing.T) {
	c := NewCache(time.Second)
	key := "1"
	nodeExp := new(Node)
	c.storage[key] = &cacheEntry{node: nodeExp, created: time.Now()}
	node, ok := c.Get(key)
	assert.Equal(t, nodeExp, node)
	assert.True(t, ok)
}

func TestCache_Get_NodeNotExists_ReturnsNilFalse(t *testing.T) {
	c := NewCache(time.Second)
	key := "1"
	node, ok := c.Get(key)
	assert.Nil(t, node)
	assert.False(t, ok)
}

func TestCache_CleanUp_ExpiredPart_RemoveExpired(t *testing.T) {
	c := NewCache(time.Minute)
	nodeExp := new(Node)
	c.storage["1"] = &cacheEntry{node: nodeExp, created: time.Now().Add(-time.Hour)}
	c.storage["2"] = &cacheEntry{node: nodeExp, created: time.Now().Add(-time.Hour)}
	c.storage["3"] = &cacheEntry{node: nodeExp, created: time.Now().Add(time.Hour)}
	c.storage["4"] = &cacheEntry{node: nodeExp, created: time.Now().Add(time.Hour)}
	c.CleanUp()
	assert.Len(t, c.storage, 2)
}
