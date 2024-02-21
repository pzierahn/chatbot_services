package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/documents"
	"log"
	"os"
)

type LegacyDocument struct {
	ID           string
	UserID       string
	CollectionID string
	Filename     string
	Path         string
}

type migrationService struct {
	db *pgxpool.Pool
}

func (mig migrationService) getAllDocuments() (docs []*LegacyDocument) {
	ctx := context.Background()
	rows, err := mig.db.Query(ctx,
		`SELECT id, user_id, collection_id, filename, path FROM documents`)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var doc LegacyDocument
		err = rows.Scan(
			&doc.ID,
			&doc.UserID,
			&doc.CollectionID,
			&doc.Filename,
			&doc.Path,
		)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}

		docs = append(docs, &doc)
	}

	return docs
}

func (mig migrationService) updateDocument(doc *LegacyDocument) {

	meta := documents.DocumentMeta{
		File: &documents.File{
			Path:     doc.Path,
			Filename: doc.Filename,
		},
	}

	ctx := context.Background()
	_, err := mig.db.Exec(ctx,
		`UPDATE documents SET metadata = $2 WHERE id = $1`,
		doc.ID, meta)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	addr := os.Getenv("CHATBOT_DB")
	db, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer db.Close()

	mig := migrationService{db: db}
	docs := mig.getAllDocuments()
	log.Printf("Found %d documents", len(docs))

	for inx, doc := range docs {
		log.Printf("Updating document[%v]: %s", inx, doc.ID)
		mig.updateDocument(doc)
	}
}
