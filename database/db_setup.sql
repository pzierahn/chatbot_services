create table if not exists collections
(
    id   uuid primary key default gen_random_uuid(),
    uid  text not null,
    name text not null
);

create table if not exists documents
(
    id         uuid primary key default gen_random_uuid(),
    uid        text not null,
    filename   text not null,
    path       text not null,
    collection uuid references collections (id) ON DELETE CASCADE
);

create table if not exists document_embeddings
(
    id        uuid primary key default gen_random_uuid(),
    source    uuid references documents (id) ON DELETE CASCADE,
    page      integer not null,
    text      text    not null,
    embedding vector(1536)
);

CREATE INDEX IF NOT EXISTS hnsw_index ON document_embeddings USING hnsw (embedding vector_ip_ops);

create table if not exists openai_usage
(
    id         uuid primary key   default gen_random_uuid(),
    created_at timestamp not null default now(),
    uid        text      not null,
    model      text      not null,
    input      int       not null,
    output     int       not null
);

create table if not exists chat_message
(
    id         uuid primary key   default gen_random_uuid(),
    uid        text      not null,
    created_at timestamp not null default now(),
    collection uuid      not null references collections (id) ON DELETE CASCADE,
    prompt     text      not null,
    completion text      not null
);

create table if not exists chat_message_source
(
    id            uuid primary key default gen_random_uuid(),
    chat          uuid not null references chat_message (id) ON DELETE CASCADE,
    document_page uuid not null references document_embeddings (id) ON DELETE CASCADE
);
