package documents

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) DeleteDocument(ctx context.Context, req *pb.Document) (*emptypb.Empty, error) {
	uid, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	_, err = service.db.Exec(ctx,
		`DELETE FROM documents WHERE id = $1 AND user_id = $2`,
		req.Id, uid)
	if err != nil {
		return nil, err
	}

	resp := service.storage.RemoveFile(bucket, []string{req.Path})
	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return &emptypb.Empty{}, nil
}
