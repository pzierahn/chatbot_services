package setup

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed db_setup.sql
var dbSetup string

// SetupTables creates all the tables in the database
func SetupTables(ctx context.Context, conn *pgxpool.Pool) error {
	_, err := conn.Exec(ctx, dbSetup)
	return err
}

func DropTables(ctx context.Context, conn *pgxpool.Pool) error {
	_, err := conn.Exec(ctx, `DROP TABLE IF EXISTS collections, documents, document_embeddings, openai_usage, chat_message, chat_message_source`)
	return err
}
