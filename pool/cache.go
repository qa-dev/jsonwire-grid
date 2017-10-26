package pool

import (
	"sync"
	"time"
)

// Cache no thread safe
type Cache struct {
	storage        map[string]*cacheEntry
	expirationTime time.Duration
	sync.RWMutex
}

type cacheEntry struct {
	node    *Node
	created time.Time
}

func NewCache(expirationTime time.Duration) *Cache {
	return &Cache{
		storage:        make(map[string]*cacheEntry),
		expirationTime: expirationTime,
	}
}

func (c *Cache) Set(key string, node *Node) {
	c.Lock()
	c.storage[key] = &cacheEntry{
		node:    node,
		created: time.Now(),
	}
	c.Unlock()
}

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

func (c *Cache) CleanUp() {
	c.Lock()
	for i, _ := range c.storage {
		if time.Since(c.storage[i].created) > c.expirationTime {
			delete(c.storage, i)
		}
	}
	c.Unlock()
}
