package gitops

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	c *github.Client
}

func NewClient(clientID, clientSecret string) *Client {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://github.com/login/oauth/access_token",
	}

	ctx := context.Background()

	oauthClient := config.Client(ctx)
	client := github.NewClient(oauthClient)

	return &Client{
		c: client,
	}
}
