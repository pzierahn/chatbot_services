package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type ScorePoints struct {
	Id     *uuid.UUID
	Source *uuid.UUID
	Page   int
	Text   string
	Score  float32
}

type SearchQuery struct {
	UserId     string
	Collection *uuid.UUID
	Embedding  []float32
	Limit      int
	Threshold  float32
}

type Page struct {
	Id    *uuid.UUID
	Page  uint32
	Text  string
	Score float32
}

type SearchResult struct {
	DocId    uuid.UUID
	Filename string
	Pages    []*Page
}

func (client *Client) Search(ctx context.Context, query SearchQuery) ([]*SearchResult, error) {
	rows, err := client.conn.Query(
		ctx,
		`select em.id, document_id, filename, page, text, (1 - (embedding <=> $1)) AS score
			from document_embeddings as em join documents as doc on doc.id = em.document_id
			where (1 - (embedding <=> $1)) >= $2 and doc.uid = $3 and doc.collection_id = $4
			ORDER BY score DESC
		 	LIMIT $5`,
		pgvector.NewVector(query.Embedding),
		query.Threshold,
		query.UserId,
		query.Collection,
		query.Limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collect := make(map[uuid.UUID]*SearchResult)
	for rows.Next() {
		var (
			id, docId uuid.UUID
			filename  string
			page      uint32
			text      string
			score     float32
		)

		if err = rows.Scan(
			&id, &docId,
			&filename,
			&page,
			&text,
			&score); err != nil {
			return nil, err
		}

		if _, ok := collect[docId]; !ok {
			collect[docId] = &SearchResult{
				DocId:    docId,
				Filename: filename,
			}
		}

		collect[docId].Pages = append(collect[docId].Pages, &Page{
			Id:    &id,
			Page:  page,
			Text:  text,
			Score: score,
		})
	}

	var results []*SearchResult
	for _, result := range collect {
		results = append(results, result)
	}

	return results, nil
}
