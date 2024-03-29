package qdrant

import (
	"context"
	"github.com/pzierahn/chatbot_services/vectordb"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

func (db *DB) Upsert(items []*vectordb.Vector) error {

	var vectors []*qdrant.PointStruct

	for _, item := range items {
		vectors = append(vectors, &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Uuid{
					Uuid: item.Id,
				},
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{
						Data: item.Vector,
					},
				},
			},
			Payload: map[string]*qdrant.Value{
				"documentId": {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.DocumentId,
					},
				},
				"collectionId": {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.CollectionId,
					},
				},
				"userId": {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.UserId,
					},
				},
				"text": {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.Text,
					},
				},
			},
		})
	}

	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"api-key",
		db.apiKey,
	)

	points := qdrant.NewPointsClient(db.conn)
	for start := 0; start < len(vectors); start += 50 {
		end := min(start+50, len(vectors))

		_, err := points.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: db.namespace,
			Points:         vectors[start:end],
		})
		if err != nil {
			return err
		}
	}

	return nil
}
