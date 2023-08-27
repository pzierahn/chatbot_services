package database

import (
	"context"
	pb "github.com/qdrant/go-client/qdrant"
	"log"
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
		Limit:          15,
		ScoreThreshold: &threshold,
		WithPayload: &pb.WithPayloadSelector{
			SelectorOptions: &pb.WithPayloadSelector_Enable{
				Enable: true,
			},
		},
	})
}

func (client *Client) DeleteFile(ctx context.Context, collection, filename string) error {
	points := pb.NewPointsClient(client.conn)

	// Delete all points with filename
	resp, err := points.Delete(ctx, &pb.DeletePoints{
		CollectionName: collection,
		Points: &pb.PointsSelector{
			PointsSelectorOneOf: &pb.PointsSelector_Filter{
				Filter: &pb.Filter{
					Should: []*pb.Condition{
						{
							ConditionOneOf: &pb.Condition_Field{
								Field: &pb.FieldCondition{
									Key: "filename",
									Match: &pb.Match{
										MatchValue: &pb.Match_Text{
											Text: filename,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}

	log.Printf("Deleted %v points\n", resp)

	return nil
}
