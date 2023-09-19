package database

import (
	"context"
	_ "embed"
)

//go:embed create_tables.sql
var tablesSQL string

// SetupTables creates all the tables in the database
func (client *Client) SetupTables(ctx context.Context) error {
	_, err := client.conn.Exec(ctx, tablesSQL)

	return err
}
