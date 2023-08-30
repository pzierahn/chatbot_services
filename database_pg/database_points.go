package database_pg

import (
	"context"
	"github.com/google/uuid"
)

type Point struct {
	Id        *uuid.UUID
	Source    uuid.UUID
	Page      int
	Text      string
	Embedding []float32
}

func (client *Client) Upsert(ctx context.Context, point Point) (uuid.UUID, error) {
	result := client.conn.QueryRow(
		ctx,
		`insert into embeddings (source, page, text, embedding)
			values ($1, $2, $3, $4) returning id`, point.Source, point.Page, point.Text, point.Embedding)

	err := result.Scan(&point.Id)
	if err != nil {
		return uuid.Nil, err
	}

	return *point.Id, nil
}

func (client *Client) SearchEmbedding(ctx context.Context, embedding []float32) ([]Point, error) {
	rows, err := client.conn.Query(
		ctx,
		`select id, source, page, text, embedding from embeddings where embedding <=> $1 < 0.8 LIMIT 15`,
		embedding)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := make([]Point, 0)
	for rows.Next() {
		point := Point{}

		err := rows.Scan(&point.Id, &point.Source, &point.Page, &point.Text, &point.Embedding)
		if err != nil {
			return nil, err
		}

		points = append(points, point)
	}

	return points, nil
}
