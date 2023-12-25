package pinecone

import (
	"context"
	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"github.com/pzierahn/chatbot_services/vectordb"
	"google.golang.org/grpc/metadata"
)

const chunkSize = 30

func (db *DB) Export(ids []string) ([]*vectordb.Vector, error) {

	if len(ids) == 0 {
		return nil, nil
	}

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", db.apiKey)

	var results []*vectordb.Vector

	for start := 0; start < len(ids); start += chunkSize {
		end := min(start+chunkSize, len(ids))

		queryResult, err := db.client.Fetch(ctx, &pinecone_grpc.FetchRequest{
			Ids:       ids[start:end],
			Namespace: "documents",
		})
		if err != nil {
			return nil, err
		}

		for _, item := range queryResult.Vectors {
			result := &vectordb.Vector{
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
	}

	return results, nil
}
