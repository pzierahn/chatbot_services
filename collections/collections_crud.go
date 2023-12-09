package collections

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/brainboost/proto"
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

	// Get all document chunk ids
	rows, err := server.db.Query(
		ctx,
		`SELECT id FROM documents WHERE collection_id = $1 AND user_id = $2`,
		collection.Id, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch document ids: %s", err)
	}

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document id: %s", err)
		}
		ids = append(ids, id)
	}

	_, err = server.db.Exec(
		ctx,
		`DELETE FROM collections WHERE id = $1 AND user_id = $2`,
		collection.Id, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete collection: %s", err)
	}

	for _, id := range ids {
		basePath := fmt.Sprintf("documents/%s/%s/%s.pdf", uid, collection.Id, id)
		log.Printf("deleting %s", basePath)

		err = server.storage.Object(basePath).Delete(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}
