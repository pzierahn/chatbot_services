package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/pzierahn/braingain/auth"
	"github.com/pzierahn/braingain/database"
	pb "github.com/pzierahn/braingain/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
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

func (server *Server) CreateCollection(ctx context.Context, collection *pb.Collection) (*emptypb.Empty, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	_, err = server.db.CreateCollection(ctx, database.Collection{
		UserId: uid.String(),
		Name:   collection.Name,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (server *Server) UpdateCollection(ctx context.Context, collection *pb.Collection) (*emptypb.Empty, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(collection.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = server.db.UpdateCollection(ctx, database.Collection{
		Id:     id,
		UserId: uid.String(),
		Name:   collection.Name,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
