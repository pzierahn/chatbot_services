package qdrant

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/search"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

func (db *Search) Search(ctx context.Context, query search.Query) (*search.Results, error) {

	embedded, err := db.embedding.CreateEmbedding(ctx, &llm.EmbeddingRequest{
		Inputs: []string{query.Query},
		Type:   llm.EmbeddingTypeQuery,
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
		Vector:         embedded.Embeddings[0],
		Limit:          uint64(query.Limit),
		Filter: &qdrant.Filter{
			Must: []*qdrant.Condition{
				{
					ConditionOneOf: &qdrant.Condition_Field{
						Field: &qdrant.FieldCondition{
							Key: search.PayloadCollectionId,
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
							Key: search.PayloadUserId,
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

	results := &search.Results{
		Results: make([]*search.Result, len(queryResult.Result)),
		Usage: search.Usage{
			ModelId: embedded.Model,
			Tokens:  embedded.Tokens,
		},
	}

	for idx, item := range queryResult.Result {
		results.Results[idx] = &search.Result{
			Id:         item.Id.GetUuid(),
			DocumentId: item.Payload[search.PayloadDocumentId].GetStringValue(),
			Text:       item.Payload[search.PayloadText].GetStringValue(),
			Position:   uint32(item.Payload[search.PayloadPosition].GetIntegerValue()),
			Score:      item.Score,
		}
	}

	return results, nil
}
