package vectordb_pinecone

import (
	"context"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
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

	var vectors []*pinecone_grpc.Vector

	for _, item := range items {
		meta := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"documentId":   {Kind: &structpb.Value_StringValue{StringValue: item.DocumentId}},
				"collectionId": {Kind: &structpb.Value_StringValue{StringValue: item.CollectionId}},
				"userId":       {Kind: &structpb.Value_StringValue{StringValue: item.UserId}},
				"filename":     {Kind: &structpb.Value_StringValue{StringValue: item.Filename}},
				"text":         {Kind: &structpb.Value_StringValue{StringValue: item.Text}},
				"page":         {Kind: &structpb.Value_NumberValue{NumberValue: float64(item.Page)}},
			},
		}

		vectors = append(vectors, &pinecone_grpc.Vector{
			Id:       item.Id,
			Values:   item.Vector,
			Metadata: meta,
		})
	}

	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"api-key",
		db.apiKey,
	)

	for start := 0; start < len(vectors); start += 50 {
		end := min(start+50, len(vectors))

		_, err := db.client.Upsert(ctx, &pinecone_grpc.UpsertRequest{
			Vectors:   vectors[start:end],
			Namespace: "documents",
		})

		if err != nil {
			return err
		}
	}

	return nil
}
