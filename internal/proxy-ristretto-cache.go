package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	ristrettoStore "github.com/eko/gocache/store/ristretto/v4"
	"log/slog"
)

// ProxyRistrettoCache implements the ProxyCache interface
type ProxyRistrettoCache struct {
	store   *cache.Cache[interface{}]
	context context.Context
}

// NewProxyRistrettoCache creates a new instance of ProxyRistrettoCache with a Redis backend
func NewProxyRistrettoCache(settings Settings) (*ProxyRistrettoCache, error) {
	slog.Info("Using Ristretto as cache store")
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: settings.Cache.Ristretto.NumCounters,
		MaxCost:     settings.Cache.Ristretto.MaxCost,
		BufferItems: settings.Cache.Ristretto.BufferItems,
	})
	if err != nil {
		panic(err)
	}

	cacheManager := cache.New[interface{}](ristrettoStore.NewRistretto(ristrettoCache))

	return &ProxyRistrettoCache{
		store:   cacheManager,
		context: context.Background(),
	}, nil
}

// Set stores a key-value pair in the Redis cache
func (p *ProxyRistrettoCache) Set(key string, value interface{}) error {
	if key == "" {
		return errors.New(ErrEmptyKey)
	}
	if value == nil {
		return errors.New("value cannot be nil")
	}
	return p.store.Set(p.context, key, value)
}

// Get retrieves a value from the Redis cache
func (p *ProxyRistrettoCache) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New(ErrEmptyKey)
	}

	value, err := p.store.Get(p.context, key)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve key '%s': %w", key, err)
	}
	return value, nil
}

// Delete removes a key-value pair from the Redis cache
func (p *ProxyRistrettoCache) Delete(key string) error {
	if key == "" {
		return errors.New(ErrEmptyKey)
	}
	return p.store.Delete(p.context, key)
}
