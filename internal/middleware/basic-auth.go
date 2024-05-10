package proxymiddleware

import (
	"basic-auth-proxy/internal/cache"
	"basic-auth-proxy/internal/settings"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

func CheckBasicAuth(cache cache.ProxyCache, settings settings.Settings) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("Checking Authorization header for Basic Auth")
			auth := r.Header.Get("Authorization")
			if auth != "" {
				isBasicAuth, credentials, err := checkForBasicAuthHeader(auth)
				if err != nil {
					slog.Debug("Invalid Authorization header")
					http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
					return
				}
				if isBasicAuth {
					if settings.Cache.Enabled {
						value, err := cache.Get(credentials)
						if err != nil || value == nil {
							slog.Debug("Cache miss")
							err := cache.Set(credentials, "true")
							if err != nil {
								slog.Debug("Error setting cache")
								http.Error(w, "Error setting cache", http.StatusInternalServerError)
								return
							}
						} else {
							slog.Debug("Cache hit")
						}
					}
				}
			} else {
				slog.Debug("No Authorization header found, skipping check")
			}
			next.ServeHTTP(w, r)
		})
	}
}

func checkForBasicAuthHeader(auth string) (bool, string, error) {
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) == 2 && parts[0] == "Basic" {
		payload, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			slog.Error("Failed to decode base64 string") // Debug log
			return false, "", errors.New("failed to decode base64 string")
		}
		return true, string(payload), nil
	}
	return false, "", nil
}
