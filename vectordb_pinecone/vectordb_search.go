package vectordb_pinecone

import (
	"context"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
	"os"
)

type SearchQuery struct {
	UserId       string
	CollectionId string
	Vector       []float32
	Limit        int
	Threshold    float32
}

func (db *DB) Search(query SearchQuery) ([]*Vector, error) {

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", os.Getenv("PINECONE_KEY"))

	queryResult, err := db.client.Query(ctx, &pinecone_grpc.QueryRequest{
		Queries: []*pinecone_grpc.QueryVector{
			{
				Values: query.Vector,
			},
		},
		Filter: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"collectionId": {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"$eq": {
									Kind: &structpb.Value_StringValue{
										StringValue: query.CollectionId,
									},
								},
							},
						},
					},
				},
				"userId": {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"$eq": {
									Kind: &structpb.Value_StringValue{
										StringValue: query.UserId,
									},
								},
							},
						},
					},
				},
			},
		},
		TopK:            200,
		IncludeValues:   false,
		IncludeMetadata: true,
		Namespace:       "documents",
	})
	if err != nil {
		return nil, err
	}

	if len(queryResult.Results) == 0 {
		return nil, nil
	}

	var results []*Vector

	for _, item := range queryResult.Results[0].Matches {
		if item.Score < query.Threshold {
			break
		}

		if len(results) >= query.Limit {
			break
		}

		doc := &Vector{
			Id:         item.Id,
			DocumentId: item.Metadata.Fields["documentId"].GetStringValue(),
			Filename:   item.Metadata.Fields["filename"].GetStringValue(),
			Page:       uint32(item.Metadata.Fields["page"].GetNumberValue()),
			Text:       item.Metadata.Fields["text"].GetStringValue(),
			Score:      item.Score,
		}

		results = append(results, doc)
	}

	return results, nil
}
