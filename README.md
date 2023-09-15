# braingain

## Setup database

### Install postgresql

```bash
docker pull ankane/pgvector
docker run -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d ankane/pgvector
```

### Start local docker container

```bash
docker compose up
```