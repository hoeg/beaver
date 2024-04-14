package gitops

import (
	"context"
	"crypto/rsa"

	"github.com/beatlabs/github-auth/app/inst"
	"github.com/google/go-github/github"
)

func NewClient(ctx context.Context, appID, instID string, privateKey *rsa.PrivateKey) (*github.Client, error) {
	conf, err := inst.NewConfig(appID, instID, privateKey)
	if err != nil {
		return nil, err
	}
	return github.NewClient(conf.Client(ctx)), nil
}
