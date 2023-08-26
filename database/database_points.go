package database

import (
	"context"
	pb "github.com/qdrant/go-client/qdrant"
)

type Payload struct {
	Uuid       string
	Collection string
	Data       []float32
	Metadata   map[string]*pb.Value
}

func (client *Client) Upsert(ctx context.Context, payload Payload) error {
	points := pb.NewPointsClient(client.conn)

	_, err := points.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: payload.Collection,
		Points: []*pb.PointStruct{
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Uuid{
						Uuid: payload.Uuid,
					},
				},
				Vectors: &pb.Vectors{
					VectorsOptions: &pb.Vectors_Vector{
						Vector: &pb.Vector{
							Data: payload.Data,
						},
					},
				},
				Payload: payload.Metadata,
			},
		},
	})

	return err
}

func (client *Client) Count(ctx context.Context, collection string) (*pb.CountResponse, error) {
	points := pb.NewPointsClient(client.conn)

	return points.Count(ctx, &pb.CountPoints{
		CollectionName: collection,
	})
}

func (client *Client) SearchEmbedding(ctx context.Context, collection string, embedding []float32) (*pb.SearchResponse, error) {
	points := pb.NewPointsClient(client.conn)

	threshold := float32(0.8)

	return points.Search(ctx, &pb.SearchPoints{
		CollectionName: collection,
		Vector:         embedding,
		Limit:          30,
		ScoreThreshold: &threshold,
		WithPayload: &pb.WithPayloadSelector{
			SelectorOptions: &pb.WithPayloadSelector_Enable{
				Enable: true,
			},
		},
	})
}
