package documents

import (
	"context"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) UpdateDocument(ctx context.Context, req *pb.Document) (*emptypb.Empty, error) {
	userID, err := service.auth.ValidateToken(ctx)
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
