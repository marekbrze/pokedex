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
	go newCache.reapLoop(interval)
	return &newCache
}

func (c *Cache) Add(key string, val []byte) {
	newCacheEntry := CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = newCacheEntry
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

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {

		c.mu.Lock()
		keysToDelete := []string{}
		for k, v := range c.cache {
			if time.Now().UTC().After(v.createdAt.Add(interval)) {
				keysToDelete = append(keysToDelete, k)
			}
		}

		if len(keysToDelete) > 0 {
			for _, k := range keysToDelete {
				delete(c.cache, k)
			}
		}

		c.mu.Unlock()
	}
}
