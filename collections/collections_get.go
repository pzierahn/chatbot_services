package collections

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
)

func (server *Service) GetCollection(ctx context.Context, col *pb.CollectionID) (*pb.Collection, error) {
	uid, err := server.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	var collection pb.Collection
	err = server.db.QueryRow(
		ctx,
		`SELECT col.id, col.name, COUNT(doc.id) AS count
			FROM collections col
			LEFT JOIN documents doc ON col.id = doc.collection_id
			WHERE col.user_id = $1 AND
			      col.id = $2
			GROUP BY col.id, col.name
			ORDER BY col.name;`,
		uid, col.Id).Scan(
		&collection.Id,
		&collection.Name,
		&collection.DocumentCount,
	)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}
