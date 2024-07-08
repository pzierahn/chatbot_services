package qdrant

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/vectordb"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

func (db *DB) Search(ctx context.Context, query vectordb.SearchQuery) (*vectordb.SearchResults, error) {

	embedding, err := db.embedding.CreateEmbedding(ctx, &llm.EmbeddingRequest{
		Inputs: []string{query.Query},
	})
	if err != nil {
		return nil, err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", db.apiKey)

	points := qdrant.NewPointsClient(db.conn)
	queryResult, err := points.Search(ctx, &qdrant.SearchPoints{
		CollectionName: db.namespace,
		WithPayload: &qdrant.WithPayloadSelector{
			SelectorOptions: &qdrant.WithPayloadSelector_Enable{
				Enable: true,
			},
		},
		ScoreThreshold: &query.Threshold,
		Vector:         embedding.Embeddings[0],
		Limit:          uint64(query.Limit),
		Filter: &qdrant.Filter{
			Must: []*qdrant.Condition{
				{
					ConditionOneOf: &qdrant.Condition_Field{
						Field: &qdrant.FieldCondition{
							Key: PayloadCollectionId,
							Match: &qdrant.Match{
								MatchValue: &qdrant.Match_Text{
									Text: query.CollectionId,
								},
							},
						},
					},
				},
				{
					ConditionOneOf: &qdrant.Condition_Field{
						Field: &qdrant.FieldCondition{
							Key: "userId",
							Match: &qdrant.Match{
								MatchValue: &qdrant.Match_Text{
									Text: query.UserId,
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(queryResult.Result) == 0 {
		return nil, nil
	}

	results := &vectordb.SearchResults{}

	for _, item := range queryResult.Result {
		fragment := &vectordb.Fragment{
			Id:           item.Id.GetUuid(),
			CollectionId: item.Payload[PayloadCollectionId].GetStringValue(),
			DocumentId:   item.Payload[PayloadDocumentId].GetStringValue(),
			UserId:       item.Payload["userId"].GetStringValue(),
			Text:         item.Payload["text"].GetStringValue(),
		}

		results.Fragments = append(results.Fragments, fragment)
		results.Scores = append(results.Scores, item.Score)
	}

	return results, nil
}
