package collections

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
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

	chunkIds, err := server.documentChunkIds(ctx, uid, collection.Id)
	if err != nil {
		return nil, err
	}

	err = server.vectorDB.Delete(chunkIds)
	if err != nil {
		return nil, err
	}

	for _, id := range docIds {
		basePath := fmt.Sprintf("documents/%s/%s/%s.pdf", uid, collection.Id, id)
		err = server.storage.Object(basePath).Delete(ctx)
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
