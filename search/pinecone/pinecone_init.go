package pinecone_search

import (
	"context"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

func (db *Search) Init() error {
	ctx := context.Background()

	// Check if the index already exists
	list, err := db.conn.ListIndexes(ctx)
	if err != nil {
		return err
	}

	for _, index := range list {
		if index.Name == db.namespace {
			return nil
		}
	}

	// Serverless index
	_, err = db.conn.CreateServerlessIndex(ctx, &pinecone.CreateServerlessIndexRequest{
		Name:               db.namespace,
		Dimension:          int32(db.dimension),
		Metric:             pinecone.Cosine,
		Cloud:              pinecone.Aws,
		Region:             "us-east-1",
		DeletionProtection: "disabled",
	})

	return err
}
