package qdrant

import (
	"context"
	"github.com/pzierahn/chatbot_services/search"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

func (db *Search) Upsert(ctx context.Context, fragments []*search.Fragment) (*search.Usage, error) {

	embedded, err := db.fastEmbedding.CreateEmbeddings(ctx, fragments)
	if err != nil {
		return nil, err
	}

	var vectors []*qdrant.PointStruct
	for _, item := range fragments {
		vector, ok := embedded.Embeddings[item.Id]
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
				search.PayloadDocumentId: {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.DocumentId,
					},
				},
				search.PayloadCollectionId: {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.CollectionId,
					},
				},
				search.PayloadUserId: {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.UserId,
					},
				},
				search.PayloadText: {
					Kind: &qdrant.Value_StringValue{
						StringValue: item.Text,
					},
				},
				search.PayloadPosition: {
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
			return nil, err
		}
	}

	return &embedded.Usage, nil
}
