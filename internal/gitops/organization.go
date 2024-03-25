package gitops

import "github.com/google/go-github/github"

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
