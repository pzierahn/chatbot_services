# Migration from 3.X.X to 4.X.X

## 1. Create a backup of the database

CockroachDB does not support the `pg_dump` command, so we need to use the `psql` command to export the data to a CSV
file.

```bash
mkdir -p backup
cd backup
psql $CHATBOT_DB -c "COPY payments TO stdout DELIMITER ',' CSV HEADER;" > payments.csv
psql $CHATBOT_DB -c "COPY collections TO stdout DELIMITER ',' CSV HEADER;" > collections.csv
psql $CHATBOT_DB -c "COPY documents TO stdout DELIMITER ',' CSV HEADER;" > documents.csv
psql $CHATBOT_DB -c "COPY document_chunks TO stdout DELIMITER ',' CSV HEADER;" > document_chunks.csv
psql $CHATBOT_DB -c "COPY model_usages TO stdout DELIMITER ',' CSV HEADER;" > model_usages.csv
psql $CHATBOT_DB -c "COPY threads TO stdout DELIMITER ',' CSV HEADER;" > threads.csv
psql $CHATBOT_DB -c "COPY thread_messages TO stdout DELIMITER ',' CSV HEADER;" > thread_messages.csv
psql $CHATBOT_DB -c "COPY thread_references TO stdout DELIMITER ',' CSV HEADER;" > thread_references.csv
psql $CHATBOT_DB -c "COPY payments TO stdout DELIMITER ',' CSV HEADER;" > payments.csv
```

## 1. Update database schema

```SQL
ALTER TABLE document_chunks
    RENAME COLUMN page TO index;

ALTER TABLE documents
    ADD COLUMN metadata jsonb NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE documents
    ADD COLUMN created_at timestamp not null default now();
```

## 2. Run the migration script

```bash
go run cmd/migration-documents/main.go
```

## 3. Delete the old columns

```SQL
ALTER TABLE documents
    DROP COLUMN filename;

ALTER TABLE documents
    DROP COLUMN path;
```