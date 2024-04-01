package gitops

import (
	"context"
	"testing"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func TestOrganizationSha(t *testing.T) {
	token := ""
	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)))

	org := NewOrganization("hoeg", "main", client)

	commitID, err := org.ReleaseBranchSha(context.Background(), "beaver")
	if err != nil {
		t.Fatal(err)
	}
	if commitID == "" {
		t.Fatal("commit id is empty")
	}
}
