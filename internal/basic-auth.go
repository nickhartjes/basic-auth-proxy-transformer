package proxy

import (
	"net/http"
)

func GetBasicAuthCredentials(r *http.Request) (username, password string, ok bool) {
	username, password, ok = r.BasicAuth()
	return
}
