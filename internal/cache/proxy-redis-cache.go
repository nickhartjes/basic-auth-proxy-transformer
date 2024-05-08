package cache

import (
	"context"
	"errors"
	"github.com/eko/gocache/lib/v4/cache"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
)

type ProxyRedisCache struct {
	store *cache.Cache[string]
}

func NewProxyRedisCache() *ProxyRedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	redisStore := redis_store.NewRedis(client)
	cacheManager := cache.New[string](redisStore)
	return &ProxyRedisCache{
		store: cacheManager,
	}
}

func (r *ProxyRedisCache) Set(key string, value interface{}) error {
	valueStr, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}
	ctx := context.Background()
	return r.store.Set(ctx, key, valueStr)
}

func (r *ProxyRedisCache) Get(key string) (interface{}, error) {
	ctx := context.Background()
	return r.store.Get(ctx, key)
}

func (r *ProxyRedisCache) Delete(key string) error {
	//TODO implement me
	panic("implement me")
}
