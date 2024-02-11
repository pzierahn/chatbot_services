package collections

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Service) GetAll(ctx context.Context, _ *emptypb.Empty) (*pb.Collections, error) {
	uid, err := server.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := server.db.Query(
		ctx,
		`SELECT col.id, col.name, COUNT(doc.id) AS count
			FROM collections col
			LEFT JOIN documents doc ON col.id = doc.collection_id
			WHERE col.user_id = $1
			GROUP BY col.id, col.name
			ORDER BY col.name;`,
		uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections pb.Collections
	for rows.Next() {
		item := new(pb.Collection)

		err = rows.Scan(
			&item.Id,
			&item.Name,
			&item.DocumentCount)
		if err != nil {
			return nil, err
		}

		collections.Items = append(collections.Items, item)
	}

	return &collections, nil
}
