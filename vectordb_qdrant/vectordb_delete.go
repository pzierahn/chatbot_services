package vectordb_qdrant

import (
	"context"
	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/metadata"
)

func (db *DB) Delete(ids []string) error {

	if len(ids) == 0 {
		return nil
	}

	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"api-key",
		db.apiKey,
	)

	var pointsIds []*qdrant.PointId
	for _, id := range ids {
		pointsIds = append(pointsIds, &qdrant.PointId{
			PointIdOptions: &qdrant.PointId_Uuid{Uuid: id},
		})
	}

	points := qdrant.NewPointsClient(db.conn)
	_, err := points.Delete(ctx, &qdrant.DeletePoints{
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Points{
				Points: &qdrant.PointsIdsList{
					Ids: make([]*qdrant.PointId, 0),
				},
			},
		},
		CollectionName: db.namespace,
	})

	return err
}
