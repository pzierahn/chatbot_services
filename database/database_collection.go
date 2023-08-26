package database

import (
	"context"
	pb "github.com/qdrant/go-client/qdrant"
	"log"
)

func (client *Client) Collections(ctx context.Context) (list []string) {
	collections := pb.NewCollectionsClient(client.conn)

	resp, err := collections.List(ctx, &pb.ListCollectionsRequest{})
	if err != nil {
		log.Fatalf("could not get collections: %v", err)
	}

	for _, collection := range resp.Collections {
		list = append(list, collection.Name)
	}

	return
}

func (client *Client) DeleteCollection(ctx context.Context, name string) error {
	collections := pb.NewCollectionsClient(client.conn)

	_, err := collections.Delete(ctx, &pb.DeleteCollection{
		CollectionName: name,
	})

	return err
}

func (client *Client) CreateCollection(ctx context.Context, name string, size uint64, distance pb.Distance) error {
	collections := pb.NewCollectionsClient(client.conn)

	_, err := collections.Create(ctx, &pb.CreateCollection{
		CollectionName: name,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     size,
					Distance: distance,
				},
			},
		},
	})

	return err
}
