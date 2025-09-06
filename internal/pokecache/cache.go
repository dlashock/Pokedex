package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type cache struct {
	entry map[string]cacheEntry
	mu    sync.Mutex
}

func NewCache(interval time.Duration) *cache {
	c := &cache{
		entry: make(map[string]cacheEntry),
	}
	return c
}

func (c *cache) Add(key string, val []byte) {
	if c == nil {
		fmt.Println("Initializing cache")
		c = NewCache(5)
	}
	var newEntry = cacheEntry{}
	newEntry.createdAt = time.Now()
	newEntry.val = val
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry[key] = newEntry
}

func (c *cache) Get(key string) ([]byte, bool) {
	if ce, exists := c.entry[key]; exists {
		c.mu.Lock()
		defer c.mu.Unlock()
		return ce.val, exists
	}
	return nil, false
}
