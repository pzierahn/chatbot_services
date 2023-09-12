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

func (client *Client) CreateCollection(ctx context.Context, coll Collection) (*uuid.UUID, error) {
	result := client.conn.QueryRow(
		ctx,
		`insert into collections (uid, name)
			values ($1, $2) returning id`,
		coll.UserId, coll.Name)

	err := result.Scan(&coll.Id)
	if err != nil {
		return nil, err
	}

	return &coll.Id, nil
}

func (client *Client) DeleteCollection(ctx context.Context, coll Collection) error {
	_, err := client.conn.Exec(
		ctx,
		`delete from collections where id = $1 and uid = $2 returning id`,
		coll.Id, coll.UserId)

	return err
}
