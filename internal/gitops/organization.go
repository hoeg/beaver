package gitops

import (
	"context"

	"github.com/google/go-github/github"
)

type Organization struct {
	name          string
	releaseBranch string
	client        *github.Client
}

func NewOrganization(name string, releaseBranch string, client *github.Client) *Organization {
	return &Organization{
		name:          name,
		releaseBranch: releaseBranch,
		client:        client,
	}
}

func (o *Organization) ReleaseBranchSha(ctx context.Context, repo string) (string, error) {
	branch, _, err := o.client.Repositories.GetBranch(ctx, o.name, repo, o.releaseBranch)
	if err != nil {
		return "", err
	}
	return *branch.Commit.SHA, nil
}
