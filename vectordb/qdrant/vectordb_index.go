package qdrant

import (
	"context"
	"github.com/pzierahn/chatbot_services/vectordb"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

func (db *DB) Upsert(ctx context.Context, fragments []*vectordb.Fragment) error {

	embeddings, err := db.createEmbeddings(ctx, fragments)
	if err != nil {
		return err
	}

	var vectors []*qdrant.PointStruct
	for _, item := range fragments {
		vector, ok := embeddings[item.Id]
		if !ok {
			continue
		}

		vectors = append(vectors, &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Uuid{
					Uuid: item.Id,
				},
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{
						Data: vector,
					},
				},
			},
			Payload: map[string]*qdrant.Value{
				PayloadDocumentId: {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.DocumentId,
					},
				},
				PayloadCollectionId: {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.CollectionId,
					},
				},
				PayloadUserId: {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.UserId,
					},
				},
				PayloadText: {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.Text,
					},
				},
				PayloadPosition: {
					Kind: &qdrant.Value_IntegerValue{
						IntegerValue: int64(item.Position),
					},
				},
			},
		})
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
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
