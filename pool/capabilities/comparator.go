package capabilities

import (
	"reflect"
	"sync"
)

type ComparatorInterface interface {
	Compare(desired Capabilities, available Capabilities) bool
	Register(Capabilities)
}

type Comparator struct {
	// полный список всевозможных зарегистрированных capabilities, для того чтобы отличать те по которым фильтруем,
	// от тех которые просто прокидываем
	registeredCaps      map[string]struct{}
	registeredCapsMutex sync.RWMutex
}

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
