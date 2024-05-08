package proxymiddleware

import (
	"basic-auth-proxy/internal/cache"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func CheckBasicAuth(cache cache.ProxyCache) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth != "" {
				isBasicAuth, credentials := checkForBasicAuthHeader(auth)
				if isBasicAuth {
					log.Printf("Basic Auth Credentials: %s\n", credentials)
					get, err := cache.Get(credentials)
					if err != nil {
						log.Printf("Cache miss")
						err := cache.Set(credentials, "true")
						if err != nil {
							log.Printf("Error setting cache: %s\n with error %s", credentials, err)
							return
						}
						get1, err1 := cache.Get(credentials)
						if err1 != nil {
							log.Printf("Error getting cache: %s\n with error %s", credentials, err)
						}
						log.Print(get1)
					}
					if get != nil {
						fmt.Print(get)
					}

					// if err != nil {
					// 	return
					// } else {
					// 	get, err := cache.Get(credentials)
					// 	if err != nil {
					// 		log.Printf("Error getting cache: %s\n with error %s", credentials, err)
					// 	}
					// 	if get != nil {
					// 		log.Printf("Cache value: %s\n", get)
					// 	}
					// }
				}
				// r.Header.Del("Authorization")
			}
			next.ServeHTTP(w, r)
		})
	}
}

func checkForBasicAuthHeader(auth string) (bool, string) {
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) == 2 && parts[0] == "Basic" {
		payload, err := base64.StdEncoding.DecodeString(parts[1])
		if err == nil {
			return true, string(payload)
		}
	}
	return false, ""
}
