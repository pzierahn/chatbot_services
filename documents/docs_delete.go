package documents

import (
	"context"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (service *Service) Delete(ctx context.Context, req *pb.Document) (*emptypb.Empty, error) {
	userId, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	ids, err := service.getChunkIds(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	err = service.db.QueryRow(ctx,
		`DELETE FROM documents
       		  WHERE id = $1 AND
					collection_id = $2 AND
					user_id = $3
 			  RETURNING path`,
		req.Id, req.CollectionId, userId).Scan(&req.Path)
	if err != nil {
		return nil, err
	}

	obj := service.storage.Object(req.Path)
	err = obj.Delete(ctx)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	err = service.vectorDB.Delete(ids)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
