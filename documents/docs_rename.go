package documents

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) Rename(ctx context.Context, req *pb.RenameDocument) (*emptypb.Empty, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	switch req.RenameTo.(type) {
	case *pb.RenameDocument_FileName:
		_, err = service.db.Exec(ctx,
			`UPDATE documents
				SET metadata = jsonb_set(metadata, '{file,filename}', to_jsonb($3::text))
				WHERE id = $1 AND 
				      user_id = $2`,
			req.Id, userId, req.GetFileName())

	case *pb.RenameDocument_WebpageTitle:
		_, err = service.db.Exec(ctx,
			`UPDATE documents
				SET metadata = jsonb_set(metadata, '{webpage,title}', to_jsonb($3::text))
				WHERE id = $1 AND 
				      user_id = $2`,
			req.Id, userId, req.GetWebpageTitle())

	default:
		err = fmt.Errorf("rename type not supported")
	}

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
