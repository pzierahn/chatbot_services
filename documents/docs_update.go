package documents

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) Update(ctx context.Context, req *pb.Document) (*emptypb.Empty, error) {
	userID, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	_, err = service.db.Exec(ctx,
		`UPDATE documents
			SET filename = $1
			WHERE id = $2 AND
			      user_id = $3 AND
			      collection_id = $4`,
		req.Filename, req.Id, userID, req.CollectionId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
