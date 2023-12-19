# Brainboost Services

This repository contains gRPC services for the Brainboost App

## Run locally

To kickstart your journey with these services, you can use the following commands:

### Set environment variables

Before you can start the server, you need to set the following environment variables:

```bash
# OpenAI API key
export OPENAI_API_KEY=""
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

# Build docker image
docker build -t brainboost .

# Use docker command (without gateway)
docker run --rm \
    -e OPENAI_API_KEY=${OPENAI_API_KEY} \
    -e SUPABASE_DB=${SUPABASE_DB} \
    -e SUPABASE_URL=${SUPABASE_URL} \
    -e SUPABASE_STORAGE_TOKEN=${SUPABASE_STORAGE_TOKEN} \
    -e SUPABASE_JWT_SECRET=${SUPABASE_JWT_SECRET} \
    -it brainboost
```

## Deploy a new service release

Prepare a new release by following these steps:

1. Update the changelog in `CHANGELOG.md`
2. Update dependencies `go get -u all`
3. Commit changes `git commit -am "Release vX.X.X"`
4. Push changes `git push`
5. Create a new git tag:
    1. `git tag -a vX.X.X -m "Release vX.X.X"`
    2. `git push origin vX.X.X`

After a new tag is pushed, the new release will be automatically build and deployed by using Google Cloud Run.

## Deploy a new gateway release

Prepare a new release by following these steps:

1. `git tag gateway/vX`
2. `git push origin gateway/vX`

After a new tag is pushed, the new release will be automatically build and deployed by using Google Cloud Run.
