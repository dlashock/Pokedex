package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entry map[string]cacheEntry
	mu    sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entry: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		for key := range c.entry {
			if time.Since(c.entry[key].createdAt) > interval {
				delete(c.entry, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) {
	var newEntry = cacheEntry{}
	newEntry.createdAt = time.Now()
	newEntry.val = val
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ce, exists := c.entry[key]; exists {
		return ce.val, exists
	}
	return nil, false
}
