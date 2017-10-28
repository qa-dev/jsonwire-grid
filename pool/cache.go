package pool

import (
	"sync"
	"time"
)

// Cache - caches nodes by session ID.
// !!! no thread safe
type Cache struct {
	storage        map[string]*cacheEntry
	expirationTime time.Duration
	sync.RWMutex
}

type cacheEntry struct {
	node    *Node
	created time.Time
}

// NewCache - constructor of Cache.
func NewCache(expirationTime time.Duration) *Cache {
	return &Cache{
		storage:        make(map[string]*cacheEntry),
		expirationTime: expirationTime,
	}
}

// Set - caches a node.
func (c *Cache) Set(key string, node *Node) {
	c.Lock()
	c.storage[key] = &cacheEntry{
		node:    node,
		created: time.Now(),
	}
	c.Unlock()
}

// Get - returns node from cache.
func (c *Cache) Get(key string) (node *Node, ok bool) {
	c.RLock()
	entry, ok := c.storage[key]
	if !ok {
		c.RUnlock()
		return nil, false
	}
	c.RUnlock()
	return entry.node, true
}

// CleanUp - removes an expired cache.
func (c *Cache) CleanUp() {
	c.Lock()
	for i := range c.storage {
		if time.Since(c.storage[i].created) > c.expirationTime {
			delete(c.storage, i)
		}
	}
	c.Unlock()
}
