package database

import (
	"context"
	"github.com/google/uuid"
)

type Collection struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
}

type CollectionInfo struct {
	ID        uuid.UUID
	Name      string
	Documents uint32
}

func (client *Client) CreateCollection(ctx context.Context, coll *Collection) (uuid.UUID, error) {
	result := client.conn.QueryRow(
		ctx,
		`insert into collections (user_id, name)
			values ($1, $2)
			returning id`,
		coll.UserID, coll.Name)

	err := result.Scan(&coll.ID)
	if err != nil {
		return uuid.Nil, err
	}

	return coll.ID, nil
}

func (client *Client) UpdateCollection(ctx context.Context, coll Collection) error {
	_, err := client.conn.Exec(
		ctx,
		`update collections set name = $3 where id = $1 and user_id = $2`,
		coll.ID, coll.UserID, coll.Name)

	return err
}

func (client *Client) DeleteCollection(ctx context.Context, coll *Collection) error {
	_, err := client.conn.Exec(
		ctx,
		`delete from collections where id = $1 and user_id = $2`,
		coll.ID, coll.UserID)

	return err
}

func (client *Client) ListCollections(ctx context.Context, uid uuid.UUID) ([]*CollectionInfo, error) {
	rows, err := client.conn.Query(
		ctx,
		`SELECT col.id, col.name, COUNT(doc.id) AS count
			FROM collections col
			LEFT JOIN documents doc ON col.id = doc.collection_id
			WHERE col.user_id = $1
			GROUP BY col.id, col.name
			ORDER BY col.name;`,
		uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collections := make([]*CollectionInfo, 0)
	for rows.Next() {
		coll := new(CollectionInfo)

		err = rows.Scan(&coll.ID, &coll.Name, &coll.Documents)
		if err != nil {
			return nil, err
		}

		collections = append(collections, coll)
	}

	return collections, nil
}
