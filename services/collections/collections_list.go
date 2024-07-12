package collections

import (
	"context"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Service) List(ctx context.Context, _ *emptypb.Empty) (*pb.CollectionList, error) {
	userId, err := server.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	collections, err := server.Database.GetCollections(ctx, userId)
	if err != nil {
		return nil, err
	}

	list := make([]*pb.Collection, len(collections))
	for idx, collection := range collections {
		list[idx] = &pb.Collection{
			Id:   collection.Id.String(),
			Name: collection.Name,
		}
	}

	return &pb.CollectionList{
		Items: list,
	}, nil
}
