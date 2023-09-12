package database

import (
	"context"
	"github.com/google/uuid"
)

type Collection struct {
	Id     uuid.UUID
	UserId string
	Name   string
}

type CollectionInfo struct {
	Id        uuid.UUID
	Name      string
	Documents uint32
}

func (client *Client) CreateCollection(ctx context.Context, coll Collection) (id *uuid.UUID, _ error) {
	result := client.conn.QueryRow(
		ctx,
		`insert into collections (uid, name)
			values ($1, $2) returning id`,
		coll.UserId, coll.Name)

	err := result.Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (client *Client) ListCollections(ctx context.Context, uid uuid.UUID) ([]*CollectionInfo, error) {
	rows, err := client.conn.Query(
		ctx,
		`select col.id, col.name, count(doc.id) as count
			from collections col join documents doc on col.id = doc.collection
			where col.uid = $1
			group by col.id, col.name`,
		uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collections := make([]*CollectionInfo, 0)
	for rows.Next() {
		coll := new(CollectionInfo)

		err := rows.Scan(&coll.Id, &coll.Name, &coll.Documents)
		if err != nil {
			return nil, err
		}

		collections = append(collections, coll)
	}

	return collections, nil
}
