package database

import (
	"context"
)

func (client *Client) SetupTables(ctx context.Context) error {
	_, err := client.conn.Exec(ctx, `
		create table if not exists collections
		(
			id   uuid primary key default gen_random_uuid(),
			uid  text not null,
			name text not null
		);
		
		create table if not exists documents
		(
			id         uuid primary key default gen_random_uuid(),
			uid  	   text not null,
			filename   text not null,
			path       text not null,
			collection uuid references collections (id)
		);
		
		create table if not exists document_embeddings
		(
			id        uuid primary key default gen_random_uuid(),
			source    uuid references documents (id),
			page      integer not null,
			text      text    not null,
			embedding vector(1536)
		);`)

	return err
}
