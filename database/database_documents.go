package database

import (
	"context"
	"github.com/google/uuid"
)

type Document struct {
	Id       uuid.UUID
	Filename string
	Tags     []string
}

func (client *Client) CreateDocument(ctx context.Context, source Document) (uuid.UUID, error) {
	result := client.conn.QueryRow(
		ctx,
		`insert into documents (filename, tags)
			values ($1, $2) returning id`, source.Filename, source.Tags)

	err := result.Scan(&source.Id)
	if err != nil {
		return uuid.Nil, err
	}

	return source.Id, nil
}

func (client *Client) ListDocuments(ctx context.Context) ([]Document, error) {
	rows, err := client.conn.Query(ctx, `select id, filename, tags from documents`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sources := make([]Document, 0)
	for rows.Next() {
		source := Document{}

		err := rows.Scan(&source.Id, &source.Filename, &source.Tags)
		if err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	return sources, nil
}

func (client *Client) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	_, err := client.conn.Exec(ctx, `delete from documents where id = $1`, id)
	return err
}

func (client *Client) GetDocument(ctx context.Context, id uuid.UUID) (*Document, error) {
	source := &Document{}

	err := client.conn.QueryRow(
		ctx,
		`select id, filename, tags from documents where id = $1`,
		id).Scan(&source.Id, &source.Filename, &source.Tags)

	if err != nil {
		return nil, err
	}

	return source, nil
}
