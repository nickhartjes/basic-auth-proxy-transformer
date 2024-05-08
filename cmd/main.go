package main

import (
	proxy "basic-auth-proxy/internal"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	settings := proxy.GetSettings()
	log.Printf("Starting proxy server on %s", settings.Port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		targetUrl, err := url.Parse(r.Header.Get("X-Target-URL"))
		if err != nil {
			http.Error(w, "Bad target URL", http.StatusBadRequest)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
