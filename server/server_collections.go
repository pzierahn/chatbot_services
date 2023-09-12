package server

import (
	"context"
	pb "github.com/pzierahn/braingain/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) ListCollections(ctx context.Context, _ *emptypb.Empty) (*pb.Collections, error) {

	collections, err := server.db.ListCollections(ctx, patrick)
	if err != nil {
		return nil, err
	}

	response := &pb.Collections{}

	for _, collection := range collections {
		response.Items = append(response.Items, &pb.Collections_Collection{
			Id:        collection.Id.String(),
			Name:      collection.Name,
			Documents: collection.Documents,
		})
	}

	return response, nil
}
