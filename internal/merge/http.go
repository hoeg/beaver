package merge

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	api "github.com/hoeg/beaver/internal/generated"
)

type API struct {
	organization string
	client       *github.Client
}

func (m *API) PostWebhook(c *gin.Context) {
	//get the request body and unmarshal it to the WebhookPayload struct

	var payload api.WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get the sha of the branch named "master" on the repository that is referenced in the payload
	sha, err := m.getMasterBranchSha(c.Request.Context(), *payload.Repository.FullName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = sha
	c.JSON(http.StatusOK, gin.H{"message": "Webhook received successfully"})
}

func (a *API) getMasterBranchSha(ctx context.Context, repository string) (string, error) {
	//get the sha of the branch named "master" on the repository that is referenced in the payload
	branch, _, err := a.client.Repositories.GetBranch(ctx, a.organization, repository, "master")
	if err != nil {
		return "", err
	}
	return *branch.Commit.SHA, nil
}
