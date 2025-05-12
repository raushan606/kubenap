package controller

import (
	"fmt"
	"sync"
	"time"
)

type ActivityCache struct {
	mu         sync.RWMutex
	timestamps map[string]time.Time // key = "namespace/service-name"
}

func NewActivityCache() *ActivityCache {
	return &ActivityCache{
		timestamps: make(map[string]time.Time),
	}
}

// Update sets the last activity timestamp for a deployment
func (c *ActivityCache) Update(namespace, name string, timestamp time.Time) {
	key := fmt.Sprintf("%s/%s", namespace, name)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timestamps[key] = timestamp
}

// Get returns the last activity timestamp, or false if not tracked
func (c *ActivityCache) Get(namespace, name string) (time.Time, bool) {
	key := fmt.Sprintf("%s/%s", namespace, name)
	c.mu.RLock()
	defer c.mu.RUnlock()
	t, ok := c.timestamps[key]
	return t, ok
}

// All returns a copy of all tracked activity timestamps
func (c *ActivityCache) All() map[string]time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]time.Time, len(c.timestamps))
	for k, v := range c.timestamps {
		out[k] = v
	}
	return out
}
