package qdrant

import (
	"context"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

func (db *Search) DeleteCollection(ctx context.Context, userId, collectionId string) error {

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		"api-key",
		db.apiKey,
	)

	points := qdrant.NewPointsClient(db.conn)
	_, err := points.Delete(ctx, &qdrant.DeletePoints{
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Filter{
				Filter: &qdrant.Filter{
					Must: []*qdrant.Condition{
						{
							ConditionOneOf: &qdrant.Condition_Field{
								Field: &qdrant.FieldCondition{
									Key: PayloadCollectionId,
									Match: &qdrant.Match{
										MatchValue: &qdrant.Match_Text{
											Text: collectionId,
										},
									},
								},
							},
						},
						{
							ConditionOneOf: &qdrant.Condition_Field{
								Field: &qdrant.FieldCondition{
									Key: PayloadUserId,
									Match: &qdrant.Match{
										MatchValue: &qdrant.Match_Text{
											Text: userId,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		CollectionName: db.namespace,
	})

	return err
}

func (db *Search) DeleteDocument(ctx context.Context, userId, documentId string) error {

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		"api-key",
		db.apiKey,
	)

	points := qdrant.NewPointsClient(db.conn)
	_, err := points.Delete(ctx, &qdrant.DeletePoints{
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Filter{
				Filter: &qdrant.Filter{
					Must: []*qdrant.Condition{
						{
							ConditionOneOf: &qdrant.Condition_Field{
								Field: &qdrant.FieldCondition{
									Key: PayloadDocumentId,
									Match: &qdrant.Match{
										MatchValue: &qdrant.Match_Text{
											Text: documentId,
										},
									},
								},
							},
						},
						{
							ConditionOneOf: &qdrant.Condition_Field{
								Field: &qdrant.FieldCondition{
									Key: PayloadUserId,
									Match: &qdrant.Match{
										MatchValue: &qdrant.Match_Text{
											Text: userId,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		CollectionName: db.namespace,
	})

	return err
}
