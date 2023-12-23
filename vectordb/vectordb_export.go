package vectordb

import (
	"context"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"google.golang.org/grpc/metadata"
	"os"
)

func (db *DB) Export(ids []string) ([]*Vector, error) {

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", os.Getenv("PINECONE_KEY"))

	queryResult, err := db.client.Fetch(ctx, &pinecone_grpc.FetchRequest{
		Ids:       ids,
		Namespace: "documents",
	})
	if err != nil {
		return nil, err
	}

	if len(queryResult.Vectors) == 0 {
		return nil, nil
	}

	var results []*Vector

	for _, item := range queryResult.Vectors {
		result := &Vector{
			Id:           item.Id,
			UserId:       item.Metadata.Fields["userId"].GetStringValue(),
			DocumentId:   item.Metadata.Fields["documentId"].GetStringValue(),
			CollectionId: item.Metadata.Fields["collectionId"].GetStringValue(),
			Filename:     item.Metadata.Fields["filename"].GetStringValue(),
			Page:         uint32(item.Metadata.Fields["page"].GetNumberValue()),
			Text:         item.Metadata.Fields["text"].GetStringValue(),
			Vector:       item.Values,
		}

		results = append(results, result)
	}

	return results, nil
}
