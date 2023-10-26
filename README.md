# Brainboost Services

This repository contains gRPC services for the Brainboost App. These services require a
running [Supabase](https://supabase.com/) instance.

## Run locally

To kickstart your journey with these services, you can use the following commands:

### Set environment variables

Before you can start the server, you need to set the following environment variables:

```bash
# OpenAI API key
export OPENAI_API_KEY=""
# Supabase database name
export SUPABASE_DB=""
# Supabase URL
export SUPABASE_URL=""
# Supabase storage token
export SUPABASE_STORAGE_TOKEN=""
# Supabase JWT secret
export SUPABASE_JWT_SECRET=""
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

### Run tests

To run the tests, run the following command:

```bash
# Tests require a running local Supabase instance
git clone https://github.com/supabase/supabase
cd supabase/docker
docker-compose up

# Run tests with
go run cmd/test/test.go
```

## Deploy a new release

Prepare a new release by following these steps:

1. Update the changelog in `CHANGELOG.md`
2. Update dependencies `go get -u all`
3. Commit changes `git commit -am "Release vX.X.X"`
4. Push changes `git push`
5. Create a new git tag:
    1. `git tag -a vX.X.X -m "Release vX.X.X"`
    2. `git push origin vX.X.X`
6. Merge `main` branch into `stable` branch

After the release is merged into the `stable` branch, the new release will be automatically deployed by using Google
Cloud Run.
