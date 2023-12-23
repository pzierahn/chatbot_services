package vectordb_qdrant

import (
	"context"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
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

	var results []*Vector

	for _, item := range queryResult.Result {
		if item.Score < query.Threshold {
			break
		}

		if len(results) >= query.Limit {
			break
		}

		doc := &Vector{
			Id:         item.Id.GetUuid(),
			DocumentId: item.Payload["documentId"].GetStringValue(),
			Filename:   item.Payload["filename"].GetStringValue(),
			Page:       uint32(item.Payload["page"].GetIntegerValue()),
			Text:       item.Payload["text"].GetStringValue(),
			Score:      item.Score,
		}

		results = append(results, doc)
	}

	return results, nil
}
