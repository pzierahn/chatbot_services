package collections

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (server *Service) Create(ctx context.Context, collection *pb.Collection) (*pb.Collection, error) {
	uid, err := server.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	result := server.db.QueryRow(
		ctx,
		`insert into collections (user_id, name)
			values ($1, $2)
			returning id`,
		uid, collection.Name)

	err = result.Scan(&collection.Id)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (server *Service) Update(ctx context.Context, collection *pb.Collection) (*pb.Collection, error) {
	uid, err := server.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	var update pb.Collection
	err = server.db.QueryRow(
		ctx,
		`UPDATE collections
			SET name = $3
			WHERE id = $1 AND user_id = $2
			RETURNING id, name`,
		collection.Id, uid, collection.Name).Scan(
		&update.Id, &update.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to update collection: %s", err)
	}

	return &update, nil
}

func (server *Service) Delete(ctx context.Context, collection *pb.Collection) (*emptypb.Empty, error) {
	uid, err := server.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	docIds, err := server.collectionDocumentIds(ctx, uid, collection.Id)
	if err != nil {
		return nil, err
	}

	chunkIds, err := server.documentChunkIds(ctx, docIds)
	if err != nil {
		return nil, err
	}

	err = server.vectorDB.Delete(chunkIds)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("documents/%s/%s", uid, collection.Id)

	iter := server.storage.Objects(ctx, &storage.Query{Prefix: basePath})
	for {
		attrs, err := iter.Next()
		if err != nil {
			break
		}

		log.Printf("Delete: %s", attrs.Name)

		err = server.storage.Object(attrs.Name).Delete(ctx)
		if err != nil {
			return nil, err
		}
	}

	_, err = server.db.Exec(
		ctx,
		`DELETE FROM collections WHERE id = $1 AND user_id = $2`,
		collection.Id, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete collection: %s", err)
	}

	return &emptypb.Empty{}, nil
}
