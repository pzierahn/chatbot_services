package qdrant

import (
	"context"
	"crypto/tls"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"os"
)

type DB struct {
	conn      *grpc.ClientConn
	apiKey    string
	namespace string
	dimension int
}

func (db *DB) Close() error {
	return db.conn.Close()
}

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

func New() (*DB, error) {

	ctx := context.Background()

	apiKey := os.Getenv("CHATBOT_QDRANT_KEY")
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", apiKey)

	target := os.Getenv("CHATBOT_QDRANT_URL")

	var opts []grpc.DialOption

	if os.Getenv("CHATBOT_QDRANT_INSECURE") == "true" {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}

	client := &DB{
		conn:      conn,
		apiKey:    apiKey,
		namespace: "documents",
		dimension: 3072,
	}
	err = client.Init()
	if err != nil {
		return nil, err
	}

	return client, nil
}
