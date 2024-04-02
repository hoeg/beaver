package gitops

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/github"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/oauth2"
)

func NewClient(ctx context.Context, appID string, privateKey []byte) *github.Client {
	ts := &tokenSource{
		appID:      appID,
		privateKey: privateKey,
		expiration: time.Minute * 10,
	}

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

type tokenSource struct {
	appID      string
	privateKey []byte
	expiration time.Duration
}

func (ts *tokenSource) Token() (*oauth2.Token, error) {
	token := jwt.New()
	_ = token.Set(jwt.IssuedAtKey, time.Now().Unix())
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(ts.expiration).Unix())
	_ = token.Set(jwt.IssuerKey, ts.appID)

	signedToken, err := jwt.Sign(token, jwa.RS256, ts.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign the JWT: %v", err)
	}

	return &oauth2.Token{
		AccessToken: string(signedToken),
	}, nil
}
