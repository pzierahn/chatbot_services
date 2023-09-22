package database

import (
	"context"
	"github.com/google/uuid"
)

type Collection struct {
	Id     uuid.UUID
	UserId uuid.UUID
	Name   string
}

type CollectionInfo struct {
	Id        uuid.UUID
	Name      string
	Documents uint32
}

func (client *Client) CreateCollection(ctx context.Context, coll *Collection) (uuid.UUID, error) {
	result := client.conn.QueryRow(
		ctx,
		`insert into collections (uid, name)
			values ($1, $2) returning id`,
		coll.UserId, coll.Name)

	err := result.Scan(&coll.Id)
	if err != nil {
		return uuid.Nil, err
	}

	return coll.Id, nil
}

func (client *Client) UpdateCollection(ctx context.Context, coll Collection) error {
	_, err := client.conn.Exec(
		ctx,
		`update collections set name = $3 where id = $1 and uid = $2`,
		coll.Id, coll.UserId, coll.Name)

	return err
}

func (client *Client) DeleteCollection(ctx context.Context, coll *Collection) error {
	_, err := client.conn.Exec(
		ctx,
		`delete from collections where id = $1 and uid = $2`,
		coll.Id, coll.UserId)

	return err
}

func (client *Client) ListCollections(ctx context.Context, uid uuid.UUID) ([]*CollectionInfo, error) {
	rows, err := client.conn.Query(
		ctx,
		`SELECT col.id, col.name, COUNT(doc.id) AS count
			FROM collections col
			LEFT JOIN documents doc ON col.id = doc.collection_id
			WHERE col.uid = $1
			GROUP BY col.id, col.name;`,
		uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collections := make([]*CollectionInfo, 0)
	for rows.Next() {
		coll := new(CollectionInfo)

		err = rows.Scan(&coll.Id, &coll.Name, &coll.Documents)
		if err != nil {
			return nil, err
		}

		collections = append(collections, coll)
	}

	return collections, nil
}
