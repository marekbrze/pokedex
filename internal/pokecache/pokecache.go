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

func NewCache(interval time.Duration) (*Cache, error) {
	return &Cache{}, nil
}
