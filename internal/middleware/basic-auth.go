package proxymiddleware

import (
	"basic-auth-proxy/internal/cache"
	"encoding/base64"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

func CheckBasicAuth(cache cache.ProxyCache) func(next http.Handler) http.Handler {
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
					value, err := cache.Get(credentials)
					if err != nil || value == nil {
						slog.Debug("Redis cache miss")
						err := cache.Set(credentials, "true")
						if err != nil {
							slog.Debug("Error setting cache")
							http.Error(w, "Error setting cache", http.StatusInternalServerError)
							return
						}
					} else {
						slog.Debug("Redis cache hit")
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
			log.Println("Failed to decode base64 string") // Debug log
			return false, "", errors.New("failed to decode base64 string")
		}
		return true, string(payload), nil
	}
	return false, "", nil
}
