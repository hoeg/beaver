# beaver

Comparing your kubernetes state with your github repositories to catch the drift.

## Kubernetes

Kubernetes resources must adhere to the following specifications:

- A label containing the repository name
- An annotation containing the commit SHA liked to the artifact

## Github App

You must install the github app with access to your github organization.