package internal

import (
	"log/slog"
)

type ProxyCache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

// SetupCache sets up the cache based on the configuration
func SetupCache(settings Settings) ProxyCache {
	var proxyCache ProxyCache
	if settings.Cache.Enabled {
		slog.Debug("Cache is enabled")
		switch settings.Cache.CacheType {
		case "redis":
			slog.Debug("Setting up Redis cache")
			proxyCache, _ = NewProxyRedisCache(settings)
		default:
			slog.Debug("Setting up Ristretto cache")
			proxyCache, _ = NewProxyRistrettoCache(settings)
		}
		slog.Debug("Cache setup complete")
	} else {
		slog.Info("Cache is disabled by configuration")
	}
	return proxyCache
}
