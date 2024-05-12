package main

import (
	"basic-auth-proxy/internal/cache"
	"basic-auth-proxy/internal/logger"
	proxymiddleware "basic-auth-proxy/internal/middleware"
	"basic-auth-proxy/internal/oauth"
	"basic-auth-proxy/internal/settings"
	"log"
	"log/slog"

	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	proxySettings := settings.GetSettings()
	logger.SetLogger(*proxySettings)
	slog.Info(fmt.Sprintf("Starting proxy server on %s", proxySettings.Port))

	proxyCache := cache.SetupCache(*proxySettings)
	oAuthConfig := oauth.GetOAuthConfig(*proxySettings)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(proxymiddleware.CheckBasicAuth(proxyCache, *proxySettings, *oAuthConfig))
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
