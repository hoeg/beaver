package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	api "github.com/hoeg/beaver/internal/generated"
	middleware "github.com/romulets/oapi-codegen/pkg/gin-middleware"

	"github.com/hoeg/beaver/internal/merge"
)

type API struct {
	merge.API
}

func New(beaverAPI api.ServerInterface) http.Handler {
	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	r := gin.Default()
	r.Use(middleware.OapiRequestValidator(swagger))
	api.RegisterHandlers(r, beaverAPI)

	return r
}
