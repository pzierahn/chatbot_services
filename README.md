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

## Build ESPv2 Google Cloud Endpoint (Not working)

[set-up-cloud-run-espv2](https://cloud.google.com/endpoints/docs/grpc/set-up-cloud-run-espv2)

```bash
# Switch to project
gcloud config set project ${PROJECT_ID}

# Enable services
gcloud services enable servicemanagement.googleapis.com
gcloud services enable servicecontrol.googleapis.com
gcloud services enable endpoints.googleapis.com

# Deploy dummy service
gcloud run deploy brainboost-gateway \
  --image="gcr.io/cloudrun/hello" \
  --allow-unauthenticated \
  --platform managed \
  --project=${PROJECT_ID}

# Deploy ESPv2 endpoint --> CONFIG
gcloud endpoints services deploy proto/api_descriptor.pb google_cloud/api_config.yaml

# Run deploy script --> IMAGE
# https://github.com/GoogleCloudPlatform/esp-v2/blob/master/docker/serverless/gcloud_build_image
./google_cloud/gcloud_build_image -s brainboost-gateway-2qkjmuus4a-ey.a.run.app \
  -c 2023-09-22r1 -p ${PROJECT_ID}

gcloud run deploy brainboost-gateway \
  --image="gcr.io/brainboost-399710/endpoints-runtime-serverless:2.45.0-brainboost-gateway-2qkjmuus4a-ey.a.run.app-2023-09-22r1" \
  --allow-unauthenticated \
  --platform managed \
  --project=${PROJECT_ID}
```
