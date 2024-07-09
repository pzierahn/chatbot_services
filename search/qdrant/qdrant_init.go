package qdrant

import (
	"context"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

func (db *DB) Init() error {
	collectionClient := qdrant.NewCollectionsClient(db.conn)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", db.apiKey)

	list, err := collectionClient.List(ctx, &qdrant.ListCollectionsRequest{})
	if err != nil {
		return err
	}

	for _, collection := range list.Collections {
		if collection.Name == db.namespace {
			//
			// Collection already exists. Nothing to do.
			//
			return nil
		}
	}

	onDisk := true
	_, err = collectionClient.Create(ctx, &qdrant.CreateCollection{
		CollectionName: db.namespace,
		VectorsConfig: &qdrant.VectorsConfig{
			Config: &qdrant.VectorsConfig_Params{
				Params: &qdrant.VectorParams{
					Size:     uint64(db.dimension),
					Distance: qdrant.Distance_Cosine,
					OnDisk:   &onDisk,
				},
			},
		},
		HnswConfig: &qdrant.HnswConfigDiff{
			OnDisk: &onDisk,
		},
	})

	return err
}
