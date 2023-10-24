# Brainboost

## Install postgresql

```bash
docker pull ankane/pgvector
docker run -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d ankane/pgvector
```

## Start docker container(s)

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

## Run Tests

Start supabase in a shell

```shell
# In one terminal
git clone https://github.com/supabase/supabase
cd supabase/docker
docker-compose up
```