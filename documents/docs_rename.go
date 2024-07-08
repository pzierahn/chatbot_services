package documents

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (service *Service) Rename(ctx context.Context, req *pb.RenameDocument) (*emptypb.Empty, error) {
	userId, err := service.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	docId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = service.Database.RenameDocument(ctx, userId, docId, req.Name)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
