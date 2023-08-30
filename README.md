# braingain

## Setup database

```sql
create table sources
(
    id       uuid primary key default gen_random_uuid(),
    filename text not null,
    tags     text[]
);

create table embeddings
(
    id        uuid primary key default gen_random_uuid(),
    source    uuid references sources (id),
    page      integer not null,
    text      text    not null,
    embedding vector(1536)
);
```