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

func processSearchQuery(rows pgx.Rows) ([]DocumentInfo, error) {
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

func (client *Client) FindDocuments(ctx context.Context, userId, like string) ([]DocumentInfo, error) {
	rows, err := client.conn.Query(ctx,
		`SELECT source, filename, collection, max(page)
		FROM documents AS doc
		    join document_embeddings AS em on doc.id = em.source
		WHERE doc.uid = $1 AND doc.filename LIKE $2
		GROUP BY source, filename, collection`,
		userId, like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return processSearchQuery(rows)
}

func (client *Client) DeleteDocument(ctx context.Context, id, uid uuid.UUID) error {
	tx, err := client.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `delete from document_embeddings where source = $1`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `delete from documents where id = $1 AND uid = $2`, id, uid)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

//func (client *Client) GetDocument(ctx context.Context, id, uid uuid.UUID) (*Document, error) {
//	source := &Document{}
//
//	err := client.conn.QueryRow(
//		ctx,
//		`select id, filename, collection from documents where id = $1 AND uid = $2`,
//		id, uid).Scan(&source.Id, &source.Filename)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return source, nil
//}

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
		    page = ANY($2)
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
