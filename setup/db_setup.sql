create table if not exists collections
(
    id      uuid primary key default gen_random_uuid(),
    user_id VARCHAR(36) not null,
    name    text        not null
);

create table if not exists documents
(
    id            uuid primary key     default gen_random_uuid(),
    user_id       VARCHAR(36) not null,
    collection_id uuid references collections (id) ON DELETE CASCADE,
    created_at    timestamp   not null default now(),
    metadata      jsonb       not null default '{}'::jsonb
);

create table if not exists document_chunks
(
    id          uuid primary key default gen_random_uuid(),
    document_id uuid references documents (id) ON DELETE CASCADE,
    text        text not null,
    index int not null
);

create table if not exists model_usages
(
    id            uuid primary key     default gen_random_uuid(),
    user_id       VARCHAR(36) not null,
    created_at    timestamp   not null default now(),
    model         text        not null,
    input_tokens  int         not null,
    output_tokens int         not null
);

create table if not exists threads
(
    id            uuid primary key     default gen_random_uuid(),
    user_id       VARCHAR(36) not null,
    created_at    timestamp   not null default now(),
    collection_id uuid        not null references collections (id) ON DELETE CASCADE
);

create table if not exists thread_messages
(
    id         uuid primary key     default gen_random_uuid(),
    user_id    VARCHAR(36) not null,
    thread_id  uuid        not null references threads (id) ON DELETE CASCADE,
    created_at timestamp   not null default now(),
    prompt     text        not null,
    completion text        not null
);

create table if not exists thread_references
(
    id                uuid primary key default gen_random_uuid(),
    user_id           VARCHAR(36) not null,
    thread_id         uuid        not null references threads (id) ON DELETE CASCADE,
    document_chunk_id uuid        not null references document_chunks (id) ON DELETE CASCADE
);

create table if not exists payments
(
    id      uuid primary key     default gen_random_uuid(),
    user_id VARCHAR(36) not null,
    date    timestamp   not null default now(),
    amount  integer     not null
);

create table if not exists crashlytics
(
    id          uuid primary key     default gen_random_uuid(),
    user_id     VARCHAR(36) not null,
    timestamp   timestamp   not null default now(),
    app_version text        not null,
    exception   text        not null,
    stack_trace text
);

create table if not exists notion_api_keys
(
    user_id VARCHAR(36) not null,
    api_key text        not null,
    primary key (user_id)
);
