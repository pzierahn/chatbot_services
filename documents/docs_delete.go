package documents

import (
	"context"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) Delete(ctx context.Context, req *pb.Document) (*emptypb.Empty, error) {
	userId, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	_, err = service.db.Exec(ctx,
		`DELETE FROM documents WHERE id = $1 AND
                            collection_id = $2 AND
                            user_id = $3`,
		req.Id, req.CollectionId, userId)
	if err != nil {
		return nil, err
	}

	obj := service.storage.Object(req.Path)
	err = obj.Delete(ctx)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
