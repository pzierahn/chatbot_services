package collections

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (server *Service) Insert(ctx context.Context, collection *pb.Collection) (*emptypb.Empty, error) {
	log.Printf("Insert: %v", collection)

	userId, err := server.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	err = server.Database.InsertCollection(ctx, &datastore.Collection{
		Id:     uuid.New(),
		UserId: userId,
		Name:   collection.Name,
	})
	if err != nil {
		log.Printf("failed to store collection: %s", err)
		return nil, fmt.Errorf("failed to store collection: %s", err)
	}

	return &emptypb.Empty{}, nil
}

func (server *Service) Update(ctx context.Context, collection *pb.Collection) (*emptypb.Empty, error) {
	log.Printf("Update: %v", collection)

	userId, err := server.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	var collectionId uuid.UUID
	if collection.Id == "" {
		collectionId = uuid.New()
	} else {
		collectionId, err = uuid.Parse(collection.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid collection id: %s", collection.Id)
		}
	}

	err = server.Database.UpdateCollection(ctx, &datastore.Collection{
		Id:     collectionId,
		UserId: userId,
		Name:   collection.Name,
	})
	if err != nil {
		log.Printf("failed to store collection: %s", err)
		return nil, fmt.Errorf("failed to store collection: %s", err)
	}

	return &emptypb.Empty{}, nil
}

func (server *Service) Delete(ctx context.Context, collection *pb.Collection) (*emptypb.Empty, error) {
	userId, err := server.Auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	collectionId, err := uuid.Parse(collection.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid collection id: %s", collection.Id)
	}

	err = server.Search.DeleteCollection(ctx, userId, collection.Id)
	if err != nil {
		return nil, err
	}

	iter := server.Storage.Objects(ctx, &storage.Query{
		Prefix: fmt.Sprintf("documents/%s/%s", userId, collection.Id),
	})
	for {
		attrs, err := iter.Next()
		if err != nil {
			break
		}

		err = server.Storage.Object(attrs.Name).Delete(ctx)
		if err != nil {
			return nil, err
		}
	}

	err = server.Database.DeleteCollection(ctx, userId, collectionId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete collection: %s", err)
	}

	return &emptypb.Empty{}, nil
}
