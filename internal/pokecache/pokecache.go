// Package pokecache is a package that helps to cashe pokeapi results
package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]CacheEntry
	mu    sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		cache: make(map[string]CacheEntry),
	}
	return &newCache
}

func (c *Cache) Add(key string, val []byte) {
	newCacheEntry := CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Lock()
	c.cache[key] = newCacheEntry
	defer c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	cacheEntry, exists := c.cache[key]
	defer c.mu.Unlock()
	if exists {
		return cacheEntry.val, true
	} else {
		return nil, false
	}
}
