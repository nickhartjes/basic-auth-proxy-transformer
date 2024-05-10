package cache

import (
	"basic-auth-proxy/internal/settings"
	"log/slog"
)

type ProxyCache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

func SetupCache(settings settings.Settings) ProxyCache {
	var proxyCache ProxyCache
	if settings.Cache.Enabled {
		switch settings.Cache.CacheType {
		case "redis":
			proxyCache, _ = NewProxyRedisCache(settings)
		default:
			proxyCache, _ = NewProxyRistrettoCache(settings)
		}
	} else {
		slog.Info("Cache is disabled by configuration")
	}
	return proxyCache
}
