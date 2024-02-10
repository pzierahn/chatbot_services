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
