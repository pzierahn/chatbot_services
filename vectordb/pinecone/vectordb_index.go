package pinecone

import (
	"context"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"github.com/pzierahn/chatbot_services/vectordb"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

func (db *DB) Upsert(items []*vectordb.Vector) error {

	var vectors []*pinecone_grpc.Vector

	for _, item := range items {
		meta := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"documentId":   {Kind: &structpb.Value_StringValue{StringValue: item.DocumentId}},
				"collectionId": {Kind: &structpb.Value_StringValue{StringValue: item.CollectionId}},
				"userId":       {Kind: &structpb.Value_StringValue{StringValue: item.UserId}},
				"text":         {Kind: &structpb.Value_StringValue{StringValue: item.Text}},
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
