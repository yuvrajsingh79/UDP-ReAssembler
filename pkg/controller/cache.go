// cache.go

package controller

import (
	"sync"
	"time"
)

// FragmentCache represents a cache for storing UDP fragments.
type FragmentCache struct {
	cache   map[string][]byte // Map of packet IDs to fragments
	mutex   sync.RWMutex
	timeout time.Duration
}

// NewFragmentCache creates a new fragment cache with a specified timeout.
func NewFragmentCache(timeout time.Duration) *FragmentCache {
	return &FragmentCache{
		cache:   make(map[string][]byte),
		timeout: timeout,
	}
}

// AddFragment adds a fragment to the cache.
func (c *FragmentCache) AddFragment(packetID string, fragment []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[packetID] = fragment

	// Set a timeout for this fragment
	go func() {
		time.Sleep(c.timeout)
		c.mutex.Lock()
		delete(c.cache, packetID)
		c.mutex.Unlock()
	}()
}

// GetFragment retrieves a fragment from the cache by packet ID.
func (c *FragmentCache) GetFragment(packetID string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	fragment, found := c.cache[packetID]
	return fragment, found
}

// RemoveFragment removes a fragment from the cache by packet ID.
func (c *FragmentCache) RemoveFragment(packetID string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, packetID)
}
