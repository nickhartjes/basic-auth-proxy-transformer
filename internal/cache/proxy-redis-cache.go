package cache

import (
	"basic-auth-proxy/internal/settings"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/eko/gocache/lib/v4/cache"
	redisStore "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
)

const ErrEmptyKey = "key cannot be empty"

// ProxyRedisCache implements the ProxyCache interface
type ProxyRedisCache struct {
	store   *cache.Cache[interface{}]
	context context.Context
}

// NewProxyRedisCache creates a new instance of ProxyRedisCache with a Redis backend
func NewProxyRedisCache(settings settings.Settings) (*ProxyRedisCache, error) {
	slog.Info("Using Redis as cache store")
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", settings.Cache.Redis.Host, settings.Cache.Redis.Port),
		Password: settings.Cache.Redis.Password,
		DB:       settings.Cache.Redis.Database,
	})

	// Check connectivity to the Redis server
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to Redis server at %s: %w", settings.Cache.Redis.Host, err)
	}
	cacheManager := cache.New[interface{}](redisStore.NewRedis(client))

	return &ProxyRedisCache{
		store:   cacheManager,
		context: context.Background(),
	}, nil
}

// Set stores a key-value pair in the Redis cache
func (p *ProxyRedisCache) Set(key string, value interface{}) error {
	if key == "" {
		return errors.New(ErrEmptyKey)
	}
	if value == nil {
		return errors.New("value cannot be nil")
	}
	return p.store.Set(p.context, key, value)
}

// Get retrieves a value from the Redis cache
func (p *ProxyRedisCache) Get(key string) (interface{}, error) {
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
func (p *ProxyRedisCache) Delete(key string) error {
	if key == "" {
		return errors.New(ErrEmptyKey)
	}
	return p.store.Delete(p.context, key)
}
