package main

import (
	"basic-auth-proxy/internal/cache"
	"basic-auth-proxy/internal/logger"
	proxymiddleware "basic-auth-proxy/internal/middleware"
	"basic-auth-proxy/internal/oauth"
	"basic-auth-proxy/internal/settings"
	"fmt"
	"log"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var proxySettings *settings.Settings

func main() {
	proxySettings = settings.GetSettings()
	logger.SetLogger(*proxySettings)
	slog.Info(fmt.Sprintf("Starting proxy server on %s", proxySettings.Port))

	proxyCache := cache.SetupCache(*proxySettings)
	oAuthConfig := oauth.GetOAuthConfig(*proxySettings)

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.RealIP, proxymiddleware.BasicAuthToOAuth(proxyCache, *proxySettings, *oAuthConfig))
	r.HandleFunc("/*", handleRequest)

	log.Fatalf("Server failed on port %s", http.ListenAndServe(fmt.Sprintf(":%s", proxySettings.Port), r))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	targetUrl, err := url.Parse(r.Header.Get(proxySettings.TargetHeaderName))
	if err != nil {
		http.Error(w, "Bad target URL", http.StatusBadRequest)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.ServeHTTP(w, r)
}
