package internal

import (
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cacheEntry map[string]CacheEntry
	interval   time.Duration
	mu         *sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cacheEntry: make(map[string]CacheEntry),
		interval:   interval,
		mu:         &sync.Mutex{},
	}

	go c.reapLoop()

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheEntry[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cacheEntry[key]
	if !exists {
		return nil, false
	}

	if time.Since(entry.createdAt) > c.interval {
		delete(c.cacheEntry, key)
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		fmt.Println("reapLoop")
		for key, entry := range c.cacheEntry {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.cacheEntry, key)
			}
		}
		c.mu.Unlock()
	}
}
