package qdrant

import (
	"context"
	"github.com/pzierahn/chatbot_services/vectordb"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

func (db *DB) Search(query vectordb.SearchQuery) ([]*vectordb.Vector, error) {

	ctx := context.Background()
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
		Vector:         query.Vector,
		Limit:          uint64(query.Limit),
		Filter: &qdrant.Filter{
			Must: []*qdrant.Condition{
				{
					ConditionOneOf: &qdrant.Condition_Field{
						Field: &qdrant.FieldCondition{
							Key: "collectionId",
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

	var results []*vectordb.Vector

	for _, item := range queryResult.Result {
		doc := &vectordb.Vector{
			Id:         item.Id.GetUuid(),
			DocumentId: item.Payload["documentId"].GetStringValue(),
			Text:       item.Payload["text"].GetStringValue(),
			Score:      item.Score,
		}

		results = append(results, doc)
	}

	return results, nil
}
