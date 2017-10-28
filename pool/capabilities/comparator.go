package capabilities

import (
	"reflect"
	"sync"
)

// ComparatorInterface - interface of capabilities comparator.
type ComparatorInterface interface {
	Compare(desired Capabilities, available Capabilities) bool
	Register(Capabilities)
}

type Comparator struct {
	// a full list of all possible registered capabilities,
	// it is necessary to distinguish those on which we filter, from those that proxy.
	registeredCaps      map[string]struct{}
	registeredCapsMutex sync.RWMutex
}

// NewComparator - constructor of comparator.
func NewComparator() *Comparator {
	return &Comparator{
		// define default capabilities for filtration
		registeredCaps: map[string]struct{}{
			"browserName":         {},
			"browserVersion":      {},
			"platformName":        {},
			"acceptInsecureCerts": {},
			"setWindowRect":       {},
		},
	}
}

// Compare - compare two set of capabilities.
func (c *Comparator) Compare(desired Capabilities, available Capabilities) bool {
	for name, currCap := range desired {
		if !c.isRegistered(name) {
			continue
		}
		if !reflect.DeepEqual(currCap, available[name]) {
			return false
		}
	}
	return true
}

// Register - registers capabilities for filtration.
func (c *Comparator) Register(caps Capabilities) {
	c.registeredCapsMutex.Lock()
	defer c.registeredCapsMutex.Unlock()
	for name := range caps {
		c.registeredCaps[name] = struct{}{}
	}
}

func (c *Comparator) isRegistered(name string) bool {
	c.registeredCapsMutex.RLock()
	defer c.registeredCapsMutex.RUnlock()
	_, ok := c.registeredCaps[name]
	return ok
}
