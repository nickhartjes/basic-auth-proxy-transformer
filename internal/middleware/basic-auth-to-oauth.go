package proxymiddleware

import (
	"basic-auth-proxy/internal/cache"
	"basic-auth-proxy/internal/settings"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"log/slog"
	"net/http"
	"strings"
)

func BasicAuthToOAuth(cache cache.ProxyCache, settings settings.Settings, oAuthConfig oauth2.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isAuthorizationHeaderPresent(r) {
				next.ServeHTTP(w, r)
				return
			}

			if isBearerToken(r) {
				next.ServeHTTP(w, r)
				return
			}

			if !isBasicAuth(r) {
				http.Error(w, "Authorization header is not Basic Auth", http.StatusBadRequest)
				return
			}

			if !isValidCredentialsFormat(r) {
				http.Error(w, "Invalid credentials format", http.StatusBadRequest)
				return
			}

			if settings.Cache.Enabled {
				handleCache(getCredentials(r), cache, oAuthConfig, r, w)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isAuthorizationHeaderPresent(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	return auth != ""
}

func isBearerToken(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	return strings.HasPrefix(auth, "Bearer ")
}

func isBasicAuth(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	isBasicAuth, _, err := checkForBasicAuthHeader(auth)
	return err == nil && isBasicAuth
}

func isValidCredentialsFormat(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	_, credentials, _ := checkForBasicAuthHeader(auth)
	creds := strings.Split(credentials, ":")
	return len(creds) == 2
}

func getCredentials(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	_, credentials, _ := checkForBasicAuthHeader(auth)
	return credentials
}

func handleCache(credentials string, cache cache.ProxyCache, oAuthConfig oauth2.Config, r *http.Request, w http.ResponseWriter) {
	value, err := cache.Get(credentials)
	if err != nil || value == nil {
		handleCacheMiss(credentials, cache, oAuthConfig, r, w)
	} else {
		handleCacheHit(value, r)
	}
}

func handleCacheMiss(credentials string, cache cache.ProxyCache, oAuthConfig oauth2.Config, r *http.Request, w http.ResponseWriter) {
	slog.Debug("Cache miss")
	creds := strings.Split(credentials, ":")
	if len(creds) != 2 {
		slog.Debug("Invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	token, tokenFetchError := oAuthConfig.PasswordCredentialsToken(context.Background(), creds[0], creds[1])
	if tokenFetchError != nil {
		fmt.Println("Error while retrieving the token:", tokenFetchError)
		return
	}

	r.Header.Del("Authorization")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	err := cache.Set(credentials, token)
	if err != nil {
		slog.Debug("Error setting cache")
		http.Error(w, "Error setting cache", http.StatusInternalServerError)
		return
	}
}

func handleCacheHit(value interface{}, r *http.Request) {
	slog.Debug("Cache hit")
	token := value.(*oauth2.Token)
	r.Header.Del("Authorization")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
}

func checkForBasicAuthHeader(auth string) (bool, string, error) {
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) == 2 && parts[0] == "Basic" {
		payload, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			slog.Error("Failed to decode base64 string")
			return false, "", errors.New("failed to decode base64 string")
		}
		return true, string(payload), nil
	}
	return false, "", nil
}
