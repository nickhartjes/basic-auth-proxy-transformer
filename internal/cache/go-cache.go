package cache

import (
	"time"

	"errors"

	"github.com/patrickmn/go-cache"
)

type GoCache struct {
	store *cache.Cache
}

func NewGoCache(defaultExpiration, cleanupInterval time.Duration) *GoCache {
	return &GoCache{
		store: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *GoCache) Set(key string, value interface{}) error {
	c.store.Set(key, value, cache.DefaultExpiration)
	return nil
}

func (c *GoCache) Get(key string) (interface{}, error) {
	value, found := c.store.Get(key)
	if !found {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (c *GoCache) Delete(key string) error {
	c.store.Delete(key)
	return nil
}
