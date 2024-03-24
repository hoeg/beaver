package app

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/hoeg/beaver/internal/generated"
	middleware "github.com/romulets/oapi-codegen/pkg/gin-middleware"

	"github.com/g4s8/go-lifecycle/pkg/adaptors"
	"github.com/gin-contrib/cors"
	"github.com/hoeg/beaver/internal/merge"
)

type API struct {
	merge.API
}

func NewAPI(beaverAPI api.ServerInterface) *adaptors.HTTPService {
	/*
		if authIssuer != "" {
			router.Use(authMiddleware(authIssuer))
		}
	*/

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	r := gin.Default()
	r.Use(middleware.OapiRequestValidator(swagger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))
	api.RegisterHandlers(r, beaverAPI)

	return adaptors.NewHTTPService(&http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: r,
	})
}
