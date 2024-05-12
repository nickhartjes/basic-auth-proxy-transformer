package oauth

import (
	"basic-auth-proxy/internal/settings"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"golang.org/x/oauth2"
)

// Update these constants with your Keycloak server and realm information
const (
	host          = "localhost"
	port          = "8090"
	tokenEndpoint = "/realms/example/protocol/openid-connect/token"
	resourceURL   = "http://localhost:8888" // Replace with the actual resource URL
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

func main() {
	// OAuth 2.0 client credentials from the Keycloak realm JSON
	clientID := "my-client"
	clientSecret := "my-client-secret"
	username := "user1"
	password := "password123"

	// Construct Keycloak token URL
	tokenURL := fmt.Sprintf("http://%s:%s"+tokenEndpoint, host, port)

	fmt.Println("Token URL:", tokenURL)
	// OAuth 2.0 configuration
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid"}, // Adjust scopes as required
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenURL,
		},
	}

	// Retrieve an access token using the Password Grant
	token, err := conf.PasswordCredentialsToken(context.Background(), username, password)
	if err != nil {
		fmt.Println("Error while retrieving the token:", err)
		return
	}

	// Display the token details
	fmt.Println("Access Token:", token.AccessToken)
	fmt.Println("Token Type:", token.TokenType)
	fmt.Println("Expiry:", token.Expiry)

	// Access a protected resource using the obtained access token
	req, _ := http.NewRequest("GET", resourceURL, nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error accessing resource:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Warn("Error closing response body:", err)
		}
	}(resp.Body)

	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return
	}
	fmt.Println("Response Data:", responseData)
}
