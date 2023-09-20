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
