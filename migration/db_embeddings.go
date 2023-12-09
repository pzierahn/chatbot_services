package migration

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	"log"
	"os"
)

type Document struct {
	Id           string
	UserId       string
	Filename     string
	Path         string
	CollectionId string
}

type DocumentEmbedding struct {
	Id         string
	Page       int
	Text       string
	Embedding  []float32
	DocumentId string
}

func ExportDocumentsMeta(ctx context.Context) []Document {
	con, err := pgxpool.New(ctx, os.Getenv("SUPABASE_DB"))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer con.Close()

	// Get all from table documents_embeddings
	rows, err := con.Query(ctx, `
		SELECT id, user_id, filename, path, collection_id
		FROM documents
	`)
	if err != nil {
		log.Fatalf("did not select: %v", err)
	}

	var docs []Document

	// Iterate over rows
	for rows.Next() {
		var document Document

		err = rows.Scan(
			&document.Id,
			&document.UserId,
			&document.Filename,
			&document.Path,
			&document.CollectionId)
		if err != nil {
			log.Fatalf("did not scan: %v", err)
		}

		docs = append(docs, document)
	}

	return docs
}

func ExportDocumentVectors(ctx context.Context, doc Document) []DocumentEmbedding {
	con, err := pgxpool.New(ctx, os.Getenv("SUPABASE_DB"))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer con.Close()

	// Get all from table documents_embeddings
	rows, err := con.Query(ctx, `
		SELECT id, page, text, embedding, document_id
		FROM document_chunks
		WHERE document_id = $1
	`, doc.Id)
	if err != nil {
		log.Fatalf("did not select: %v", err)
	}

	var embeddings []DocumentEmbedding

	// Iterate over rows
	for rows.Next() {
		var embedding DocumentEmbedding
		var vector pgvector.Vector

		err = rows.Scan(
			&embedding.Id,
			&embedding.Page,
			&embedding.Text,
			&vector,
			&embedding.DocumentId)
		if err != nil {
			log.Fatalf("did not scan: %v", err)
		}

		embedding.Embedding = vector.Slice()

		embeddings = append(embeddings, embedding)
	}

	return embeddings
}
