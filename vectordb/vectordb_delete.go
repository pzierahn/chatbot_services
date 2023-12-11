package vectordb

import (
	"context"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"google.golang.org/grpc/metadata"
)

func (db *DB) Delete(ids []string) error {

	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"api-key",
		db.apiKey,
	)

	_, err := db.client.Delete(ctx, &pinecone_grpc.DeleteRequest{
		Ids:       ids,
		DeleteAll: false,
		Namespace: "documents",
	})

	return err
}

func (db *DB) DeleteAll() error {

	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"api-key",
		db.apiKey,
	)

	_, err := db.client.Delete(ctx, &pinecone_grpc.DeleteRequest{
		DeleteAll: true,
		Namespace: "documents",
	})

	return err
}
