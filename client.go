package pcloud

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type ClientConfig struct {
	ClientID     string
	ClientSecret string
}

type pCloudClient struct {
	Client *http.Client
}

// NewClient create new pCloudClient
func NewClient(c ClientConfig) *pCloudClient {
	ctx := context.Background()

	conf := oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://my.pcloud.com/oauth2/authorize",
			TokenURL: "https://api.pcloud.com/oauth2_token",
		},
	}

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	return &pCloudClient{
		Client: client,
	}
}
