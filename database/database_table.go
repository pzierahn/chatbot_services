package database

import (
	"context"
)

func (client *Client) CreateExtension(ctx context.Context) error {
	_, err := client.conn.Exec(ctx, `CREATE EXTENSION if not exists vector;`)
	return err
}

func (client *Client) CreateTables(ctx context.Context) error {
	_, err := client.conn.Exec(
		ctx,
		`create table documents
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
			);`)

	return err
}

func (client *Client) DropTables(ctx context.Context) error {
	_, err := client.conn.Exec(
		ctx,
		`drop table if exists document_embeddings; drop table if exists documents;`)

	return err
}
