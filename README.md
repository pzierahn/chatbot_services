# Brainboost Services

This repository contains gRPC services for the Brainboost App

## Run locally

### Set up environment

A Google [Firebase](https://firebase.google.com/) Service account is needed for Authentication and Storage.
Store the service account credentials in `./service_account.json` in the main folder.

You also need to set the following environment variables:

```bash
# OpenAI API key
export OPENAI_API_KEY=""

# QDRANT API key
export CHATBOT_QDRANT_KEY=""

# QDRANT API URL
export CHATBOT_QDRANT_URL=""

# Postgres database connection string
export CHATBOT_DB=""

# AWS Bedrock Credentials
export AWS_ACCESS_KEY_ID=""
export AWS_SECRET_ACCESS_KEY=""
```

### Start the server

To start the server, run the following command:

```bash
go run cmd/server/server.go
```

Or use the following command to start the server with an envoy proxy:

```bash
# Use docker-compose to start service and gateway
docker compose up
```

## Deploy a new service release

Prepare a new release by following these steps:

1. Update the changelog in `CHANGELOG.md`
2. Update dependencies `go get -u all`
3. Commit changes `git commit -am "Release vX.X.X"`
4. Push changes `git push`
5. Create a new git tag:
    1. `git tag vX.X.X`
    2. `git push origin vX.X.X`

After a new tag is pushed, the new release will be automatically build and deployed by using Google Cloud Run.

## Deploy a new gateway release

To use gRPC services in browser a gRPC-Web translator is needed. These Proxies are documented in `envoy/`.

To prepare a new gateway release run the following steps:

1. `git tag gateway/vX`
2. `git push origin gateway/vX`

After a new tag is pushed, the new release will be automatically build and deployed by using Google Cloud Run.

## Local development

For local testing run the following command:

```shell
# Start the server
CHATBOT_TEST=true \
CHATBOT_QDRANT_INSECURE=true \
PORT=8869 \
CHATBOT_DB=postgresql://root@127.0.0.1:26257/defaultdb?sslmode=disable \
CHATBOT_QDRANT_URL=localhost:6334 \
go run cmd/server/server.go
```