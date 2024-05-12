package oauth

import (
	"basic-auth-proxy/internal/settings"
	"fmt"
	"golang.org/x/oauth2"
)

func GetOAuthConfig(settings settings.Settings) *oauth2.Config {
	tokenURL := fmt.Sprintf("http://%s:%d:%s", settings.OAuth2.Host, settings.OAuth2.Port, settings.OAuth2.TokenEndpoint)
	conf := &oauth2.Config{
		ClientID:     settings.OAuth2.ClientID,
		ClientSecret: settings.OAuth2.ClientSecret,
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenURL,
		},
	}
	return conf
}
