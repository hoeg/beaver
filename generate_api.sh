#!/bin/bash

# go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
oapi-codegen -package api -generate types api/openapi.yaml > internal/generated/types.go
oapi-codegen -package api -generate spec api/openapi.yaml > internal/generated/spec.go 
oapi-codegen -package api -generate gin api/openapi.yaml > internal/generated/api.go