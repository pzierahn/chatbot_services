package chat

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
)

func (service *Service) GetThreads(ctx context.Context, collection *pb.Collection) (*pb.ThreadIDs, error) {
	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(
		ctx,
		`SELECT id
			FROM threads
			WHERE user_id = $1 AND collection_id = $2
			ORDER BY created_at DESC`,
		userId, collection.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := &pb.ThreadIDs{}

	for rows.Next() {
		var threadId string
		err = rows.Scan(&threadId)
		if err != nil {
			return nil, err
		}

		list.Ids = append(list.Ids, threadId)
	}

	return list, nil
}
