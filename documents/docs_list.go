package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
)

func (service *Service) List(ctx context.Context, req *pb.DocumentFilter) (*pb.DocumentList, error) {

	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := service.db.Query(ctx,
		`SELECT id, title
		FROM documents
		WHERE
		    user_id = $1 AND
		    collection_id = $2::uuid AND
		    title LIKE $3`,
		userId, req.CollectionId, "%"+req.Query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := &pb.DocumentList{
		Items: make(map[string]string),
	}

	for rows.Next() {
		var (
			docId string
			title string
		)

		err = rows.Scan(
			&docId,
			&title,
		)
		if err != nil {
			return nil, err
		}

		documents.Items[docId] = title
	}

	return documents, nil
}
