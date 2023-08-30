package database_pg

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Client struct {
	conn *pgxpool.Pool
}

func (client *Client) Close() {
	client.conn.Close()
}

func Connect(ctx context.Context, addr string) (*Client, error) {
	conn, err := pgxpool.Connect(ctx, addr)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{conn: conn}, nil
}
