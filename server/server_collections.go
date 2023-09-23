package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/database"
	pb "github.com/pzierahn/brainboost/proto"
	supastorage "github.com/supabase-community/storage-go"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (server *Server) GetCollections(ctx context.Context, _ *emptypb.Empty) (*pb.Collections, error) {

	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	collections, err := server.db.ListCollections(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := &pb.Collections{}

	for _, collection := range collections {
		response.Items = append(response.Items, &pb.Collections_Collection{
			Id:        collection.ID.String(),
			Name:      collection.Name,
			Documents: collection.Documents,
		})
	}

	return response, nil
}

func (server *Server) CreateCollection(ctx context.Context, collection *pb.Collection) (*emptypb.Empty, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	_, err = server.db.CreateCollection(ctx, &database.Collection{
		UserID: uid,
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
		ID:     id,
		UserID: uid,
		Name:   collection.Name,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (server *Server) DeleteCollection(ctx context.Context, collection *pb.Collection) (*emptypb.Empty, error) {
	uid, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(collection.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = server.db.DeleteCollection(ctx, &database.Collection{
		ID:     id,
		UserID: uid,
	})
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("%s/%s", uid, id)

	var paths []string
	fileObjs := server.storage.ListFiles(bucket, basePath, supastorage.FileSearchOptions{})
	for _, file := range fileObjs {
		paths = append(paths, basePath+"/"+file.Name)
	}

	resp := server.storage.RemoveFile(bucket, paths)
	if resp.Error != "" {
		return nil, fmt.Errorf("failed to delete file: %s", resp.Error)
	}

	return &emptypb.Empty{}, nil
}
