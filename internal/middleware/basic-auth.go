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

func CheckBasicAuth(cache cache.ProxyCache, settings settings.Settings, oAuthConfig oauth2.Config) func(next http.Handler) http.Handler {
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

							// credentials split on : the first values is username the second is password
							creds := strings.Split(credentials, ":")
							if len(creds) != 2 {
								slog.Debug("Invalid credentials")
								http.Error(w, "Invalid credentials", http.StatusBadRequest)
								return
							}

							// Retrieve an access token using the Password Grant
							token, tokenFetchError := oAuthConfig.PasswordCredentialsToken(context.Background(), creds[0], creds[1])
							if tokenFetchError != nil {
								fmt.Println("Error while retrieving the token:", tokenFetchError)
								return
							} else {
								// add the token to the header as a bearer token
								r.Header.Del("Authorization")
								r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
								// remove the basic auth header

								// Display the token details
								fmt.Println("Access Token:", token.AccessToken)
								fmt.Println("Token Type:", token.TokenType)
								fmt.Println("Expiry:", token.Expiry)

								error11 := cache.Set(credentials, token)
								r.Header.Del("Authorization")
								r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
								if error11 != nil {
									slog.Debug("Error setting cache")
									http.Error(w, "Error setting cache", http.StatusInternalServerError)
									return
								}
							}

						} else {
							slog.Debug("Cache hit")

							// convert the value to a token
							value := value.(*oauth2.Token)
							fmt.Println("Cache Access Token:", value.AccessToken)
							fmt.Println("Cache Token Type:", value.TokenType)
							fmt.Println("Cache Expiry:", value.Expiry)
							r.Header.Del("Authorization")
							r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", value.AccessToken))
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
