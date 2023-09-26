package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Database interface {
	Close()
}

type Client struct {
	conn *pgxpool.Pool
}

func (client *Client) Close() {
	client.conn.Close()
}

func Connect(ctx context.Context, addr string) (*Client, error) {
	conn, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{conn: conn}, nil
}
