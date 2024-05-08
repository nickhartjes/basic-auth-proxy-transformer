package main

import (
	"basic-auth-proxy/internal/cache"
	proxymiddleware "basic-auth-proxy/internal/middleware"
	settings "basic-auth-proxy/internal/settings"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
)

func main() {
	proxySettings := settings.GetSettings()
	log.Printf("Starting proxy server on %s", proxySettings.Port)

	var mycache cache.ProxyCache
	mycache = cache.NewGoCache(5, 10)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(proxymiddleware.CheckBasicAuth(mycache))
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		targetUrl, err := url.Parse(r.Header.Get("X-Target-URL"))
		if err != nil {
			http.Error(w, "Bad target URL", http.StatusBadRequest)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", proxySettings.Port), r))
}
