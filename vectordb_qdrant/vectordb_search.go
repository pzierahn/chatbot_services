package vectordb_qdrant

import (
	"context"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

type SearchResults struct {
	Items []*Document
}

type Document struct {
	Id         string
	DocumentId string
	Filename   string
	Page       uint32
	Content    string
	Score      float32
}

type SearchQuery struct {
	UserId       string
	CollectionId string
	Vector       []float32
	Limit        int
	Threshold    float32
}

func (db *DB) Search(query SearchQuery) (*SearchResults, error) {

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
		return &SearchResults{}, nil
	}

	results := &SearchResults{}

	for _, item := range queryResult.Result {
		if item.Score < query.Threshold {
			break
		}

		if len(results.Items) >= query.Limit {
			break
		}

		doc := &Document{
			Id:         item.Id.GetUuid(),
			DocumentId: item.Payload["documentId"].GetStringValue(),
			Filename:   item.Payload["filename"].GetStringValue(),
			Page:       uint32(item.Payload["page"].GetIntegerValue()),
			Content:    item.Payload["text"].GetStringValue(),
			Score:      item.Score,
		}

		results.Items = append(results.Items, doc)
	}

	return results, nil
}
