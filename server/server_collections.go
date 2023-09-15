package server

import (
	"context"
	"github.com/pzierahn/braingain/auth"
	pb "github.com/pzierahn/braingain/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) GetCollections(ctx context.Context, _ *emptypb.Empty) (*pb.Collections, error) {

	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	collections, err := server.db.ListCollections(ctx, *uid)
	if err != nil {
		return nil, err
	}

	response := &pb.Collections{}

	for _, collection := range collections {
		response.Items = append(response.Items, &pb.Collections_Collection{
			Id:        collection.Id.String(),
			Name:      collection.Name,
			Documents: uint32(collection.Documents),
		})
	}

	return response, nil
}
