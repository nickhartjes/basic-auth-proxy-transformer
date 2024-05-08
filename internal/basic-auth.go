package proxy

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

func CheckBasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			isBasicAuth, credentials := checkForBasicAuthHeader(auth)
			if isBasicAuth {
				log.Printf("Basic Auth Credentials: %s\n", credentials)
			}
		}
		next.ServeHTTP(w, r)
	})
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
