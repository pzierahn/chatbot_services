package vectordb

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
)

type DB struct {
	conn   *grpc.ClientConn
	apiKey string
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func New() (*DB, error) {

	ctx := context.Background()

	apiKey := os.Getenv("QDRANT_KEY")
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", apiKey)

	target := os.Getenv("QDRANT_URL")

	log.Printf("connecting to %v", target)

	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	)
	if err != nil {
		return nil, err
	}

	return &DB{
		conn:   conn,
		apiKey: apiKey,
	}, nil
}
