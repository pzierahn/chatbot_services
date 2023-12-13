# Migration Guide From Supabase to Firebase

## Database

To run database migrations, you can use the following command:

```shell
docker run --rm -it \
  -e SUPABASE_DB=$SUPABASE_DB \
  -e BRAINBOOST_COCKROACH_DB=$BRAINBOOST_COCKROACH_DB \
  postgres /bin/bash

pg_dump -Fc -v -d $SUPABASE_DB -f brainboost.dump
pg_restore -v -d $AWS_BRAINBOOST_DB brainboost.dump

# Export the table data to a CSV file
psql $SUPABASE_DB -c "\COPY xxxxx_document_embeddings_copy TO 'xxxxx_document_embeddings_copy.csv' WITH CSV HEADER;"

# Import the table data from a CSV file
psql $BRAINBOOST_COCKROACH_DB -c "\COPY document_embeddings FROM 'xxxxx_document_embeddings_copy.csv' WITH CSV HEADER;"
```

```sql
create table if not exists xxxxx_document_embeddings_copy
(
    id          uuid primary key default gen_random_uuid(),
    document_id uuid    not null,
    page        integer not null,
    text        text    not null
);

INSERT INTO xxxxx_document_embeddings_copy
SELECT id, document_id, page, text FROM document_embeddings;

INSERT INTO document_chunks
SELECT * FROM document_embeddings;


create table if not exists document_embeddings
(
    id          uuid primary key default gen_random_uuid(),
    document_id uuid    not null,
    page        integer not null,
    text        text    not null
);

```
