package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val []byte
	createdAt time.Time
}

type Cache struct {
	mu sync.RWMutex
	entries map[string]cacheEntry
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
		interval: interval,
	}
	go c.readLoop()

	return c

}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()

	c.entries[key] = cacheEntry{
		val: val,
		createdAt: time.Now(),
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	data, ok := c.entries[key]
    c.mu.RUnlock()
	return data.val, ok
}

func (c *Cache) readLoop() {
	ticker := time.NewTicker(c.interval)

    for range ticker.C {
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.entries, key)
			}
		}
	}
}