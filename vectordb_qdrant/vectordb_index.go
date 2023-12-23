package vectordb_qdrant

import (
	"context"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

type Vector struct {
	Id           string
	DocumentId   string
	UserId       string
	CollectionId string
	Filename     string
	Text         string
	Page         uint32
	Vector       []float32
	Score        float32
}

func (db *DB) Upsert(items []*Vector) error {

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
						Data: []float32{1, 2, 3},
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
				"filename": {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.Filename,
					},
				},
				"text": {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.Text,
					},
				},
				"page": {
					Kind: &qdrant.Value_IntegerValue{
						IntegerValue: int64(item.Page),
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
