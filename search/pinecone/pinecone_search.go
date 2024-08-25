package pinecone_search

import (
	"context"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/search"
	"google.golang.org/protobuf/types/known/structpb"
)

func (db *Search) Search(ctx context.Context, query search.Query) (*search.Results, error) {

	embedded, err := db.embedding.CreateEmbedding(ctx, &llm.EmbeddingRequest{
		Inputs: []string{query.Query},
		Type:   llm.EmbeddingTypeQuery,
	})
	if err != nil {
		return nil, err
	}

	idxConnection, err := db.getIndexConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = idxConnection.Close() }()

	filter, err := structpb.NewStruct(map[string]any{
		search.PayloadCollectionId: query.CollectionId,
		search.PayloadUserId:       query.UserId,
	})
	if err != nil {
		return nil, err
	}

	vectors, err := idxConnection.QueryByVectorValues(ctx, &pinecone.QueryByVectorValuesRequest{
		Vector:          embedded.Embeddings[0],
		TopK:            query.Limit,
		MetadataFilter:  filter,
		IncludeValues:   false,
		IncludeMetadata: true,
	})
	if err != nil {
		return nil, err
	}

	var results []*search.Result
	for _, match := range vectors.Matches {
		if match.Vector == nil || match.Score < query.Threshold {
			continue
		}

		vector := match.Vector
		metadata := vector.Metadata.AsMap()

		fragmentId, ok := metadata[search.PayloadFragmentId].(string)
		if !ok {
			continue
		}

		text, ok := metadata[search.PayloadText].(string)
		if !ok {
			continue
		}

		documentId, ok := metadata[search.PayloadDocumentId].(string)
		if !ok {
			continue
		}

		posFloat, ok := metadata[search.PayloadPosition].(float64)
		if !ok {
			continue
		}

		position := uint32(posFloat)

		results = append(results, &search.Result{
			Id:         fragmentId,
			Text:       text,
			DocumentId: documentId,
			Position:   position,
			Score:      match.Score,
		})
	}

	return &search.Results{
		Usage: search.Usage{
			ModelId: embedded.Model,
			Tokens:  embedded.Tokens,
		},
		Results: results,
	}, nil
}
