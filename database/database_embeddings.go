package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type Point struct {
	Id        *uuid.UUID
	Source    uuid.UUID
	Page      int
	Text      string
	Embedding []float32
}

type ScorePoints struct {
	Id     uuid.UUID
	Source uuid.UUID
	Page   int
	Text   string
	Score  float32
}

type SearchQuery struct {
	Embedding []float32
	Limit     int
	Threshold float32
}

func (client *Client) Upsert(ctx context.Context, point Point) (uuid.UUID, error) {
	result := client.conn.QueryRow(
		ctx,
		`insert into document_embeddings (source, page, text, embedding)
			values ($1, $2, $3, $4) returning id`, point.Source, point.Page, point.Text, pgvector.NewVector(point.Embedding))

	err := result.Scan(&point.Id)
	if err != nil {
		return uuid.Nil, err
	}

	return *point.Id, nil
}

func (client *Client) SearchEmbedding(ctx context.Context, query SearchQuery) ([]ScorePoints, error) {
	rows, err := client.conn.Query(
		ctx,
		`select id, source, page, text, (1 - (embedding <=> $1)) AS score
			from document_embeddings
			where (1 - (embedding <=> $1)) >= $2
			ORDER BY score DESC
		 	LIMIT $3`,
		pgvector.NewVector(query.Embedding),
		query.Threshold,
		query.Limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := make([]ScorePoints, 0)
	for rows.Next() {
		point := ScorePoints{}

		err = rows.Scan(&point.Id, &point.Source, &point.Page, &point.Text, &point.Score)
		if err != nil {
			return nil, err
		}

		points = append(points, point)
	}

	return points, nil
}
