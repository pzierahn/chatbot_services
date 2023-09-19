package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

type PageEmbedding struct {
	Page      int
	Text      string
	Embedding []float32
}

type Document struct {
	Id         uuid.UUID
	UserId     string
	Collection uuid.UUID
	Filename   string
	Path       string
	Pages      []*PageEmbedding
}

type DocumentInfo struct {
	Id         uuid.UUID
	Collection uuid.UUID
	Filename   string
	Pages      uint32
}

type DocumentQuery struct {
	UserId     string
	Collection *uuid.UUID
	Query      string
}

func (client *Client) UpsertDocument(ctx context.Context, doc Document) (*uuid.UUID, error) {
	tx, err := client.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	result := client.conn.QueryRow(
		ctx,
		`insert into documents (uid, filename, path, collection)
			values ($1, $2, $3, $4) returning id`,
		doc.UserId,
		doc.Filename,
		doc.Path,
		doc.Collection)
	err = result.Scan(&doc.Id)
	if err != nil {
		return nil, err
	}

	for _, page := range doc.Pages {
		_, err = tx.Exec(ctx,
			`insert into document_embeddings (source, page, text, embedding)
				values ($1, $2, $3, $4)`,
			doc.Id,
			page.Page,
			page.Text,
			pgvector.NewVector(page.Embedding))
		if err != nil {
			return nil, err
		}
	}

	return &doc.Id, tx.Commit(ctx)
}

func (client *Client) FindDocuments(ctx context.Context, query DocumentQuery) ([]DocumentInfo, error) {
	rows, err := client.conn.Query(ctx,
		`SELECT source, filename, collection, max(page)
		FROM documents AS doc
		    join document_embeddings AS em on doc.id = em.source
		WHERE
		    doc.uid = $1 AND
		    ($2::uuid is null OR doc.collection = $2::uuid) AND
		    doc.filename LIKE $3
		GROUP BY source, filename, collection`,
		query.UserId, query.Collection, query.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sources := make([]DocumentInfo, 0)

	for rows.Next() {
		source := DocumentInfo{}

		err := rows.Scan(
			&source.Id,
			&source.Filename,
			&source.Collection,
			&source.Pages)
		if err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	return sources, nil
}

func (client *Client) DeleteDocument(ctx context.Context, id uuid.UUID, uid string) error {
	_, err := client.conn.Exec(ctx, `delete from documents where id = $1 and uid = $2`, id, uid)
	return err
}

type PageContentQuery struct {
	Id     uuid.UUID
	UserId string
	Pages  []uint32
}

type PageContent struct {
	Id       uuid.UUID
	Filename string
	Page     uint32
	Text     string
}

func (client *Client) GetPageContent(ctx context.Context, query PageContentQuery) ([]*PageContent, error) {
	rows, err := client.conn.Query(ctx,
		`select dm.id, filename, page, text
		from document_embeddings as dm, documents as doc
		where
		    source = $1 and
		    doc.id = dm.source and
		    uid = $2 and
		    page = ANY($3)
		order by filename, page`, query.Id, query.UserId, query.Pages)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sources := make([]*PageContent, 0)
	for rows.Next() {
		source := &PageContent{}
		err = rows.Scan(
			&source.Id,
			&source.Filename,
			&source.Page,
			&source.Text)
		if err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	return sources, nil
}
