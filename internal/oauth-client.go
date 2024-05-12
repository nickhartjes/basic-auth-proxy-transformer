package internal

import (
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/url"
)

// GetOAuthConfig creates an OAuth2 configuration from the settings
func GetOAuthConfig(settings Settings) *oauth2.Config {
	// Create the token URL
	tokenURL := fmt.Sprintf("%s:%d:%s", settings.OAuth2.Host, settings.OAuth2.Port, settings.OAuth2.TokenEndpoint)

	// Check if the URL is valid
	_, err := url.Parse(tokenURL)
	if err != nil {
		log.Fatalf("invalid token URL: %v", err)
	}

	// Create the OAuth2 configuration
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
