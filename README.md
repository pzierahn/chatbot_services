# braingain

## Setup database

### Install postgresql

```bash
docker pull ankane/pgvector
docker run -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d ankane/pgvector
```

### Create database tables

```sql
create table documents
(
    id       uuid primary key default gen_random_uuid(),
    filename text not null,
    tags     text[]
);

create table document_embeddings
(
    id        uuid primary key default gen_random_uuid(),
    source    uuid references documents (id),
    page      integer not null,
    text      text    not null,
    embedding vector(1536)
);
```

### Delete database tables

```sql
DROP TABLE document_embeddings;
DROP TABLE documents;
```