#!/bin/bash

#build Docker image
docker build -t beaver:test .

# Create a Kind cluster
kind create cluster --name my-cluster

# Load the Docker image into the Kind cluster
kind load docker-image beaver:test --name my-cluster

# apply all the resources in this folder except test-deployment.yaml
kubectl apply -f test/clusterrole.yaml
kubectl apply -f test/serviceaccount.yaml
kubectl apply -f test/clusterrolebinding.yaml
kubectl apply -f test/deployment.yaml

# wait for the pods to be ready
kubectl wait --for=condition=ready --timeout=5m pod --all

