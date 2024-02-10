package chat

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) DeleteThread(ctx context.Context, req *pb.ThreadID) (*emptypb.Empty, error) {
	userId, err := service.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	_, err = service.db.Exec(ctx,
		`DELETE FROM threads
			WHERE user_id = $1 AND
			      id = $2`,
		userId, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
