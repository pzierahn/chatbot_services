package pinecone

import (
	"context"
	"crypto/tls"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
)

type DB struct {
	client pinecone_grpc.VectorServiceClient
	conn   *grpc.ClientConn
	apiKey string
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func New() (*DB, error) {
	config := &tls.Config{}

	apiKey := os.Getenv("PINECONE_KEY")
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", apiKey)
	target := os.Getenv("PINECONE_URL")

	log.Printf("connecting to %v", target)

	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithAuthority(target),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	client := pinecone_grpc.NewVectorServiceClient(conn)

	return &DB{
		client: client,
		conn:   conn,
		apiKey: apiKey,
	}, nil
}
