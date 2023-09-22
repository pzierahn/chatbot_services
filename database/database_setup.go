package database

import (
	"context"
	_ "embed"
)

//go:embed db_setup.sql
var dbSetup string

// SetupTables creates all the tables in the database
func (client *Client) SetupTables(ctx context.Context) error {
	_, err := client.conn.Exec(ctx, dbSetup)
	return err
}

func (client *Client) DropTables(ctx context.Context) error {
	_, err := client.conn.Exec(ctx, `DROP TABLE IF EXISTS collections, documents, document_embeddings, openai_usage, chat_message, chat_message_source`)
	return err
}
