package merge

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/hoeg/beaver/internal/generated"
)

type API struct {
}

func (m *API) PostWebhook(c *gin.Context) {
	//get the request body and unmarshal it to the WebhookPayload struct

	var payload api.WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// process the payload and perform necessary actions
	// ...
	c.JSON(http.StatusOK, gin.H{"message": "Webhook received successfully"})

}
