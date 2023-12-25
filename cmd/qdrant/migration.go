package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/vectordb/pinecone"
	vectordb_qdrant2 "github.com/pzierahn/chatbot_services/vectordb/qdrant"
	"log"
	"os"
)

func getDocs(ctx context.Context, db *pgxpool.Pool) []string {
	rows, err := db.Query(ctx, `SELECT id FROM documents`)
	if err != nil {
		log.Fatal(err)
	}

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}

	return ids
}

func getChunks(ctx context.Context, db *pgxpool.Pool, docId string) []string {
	rows, err := db.Query(ctx, `SELECT id FROM document_chunks WHERE document_id = $1`, docId)
	if err != nil {
		log.Fatal(err)
	}

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}

	return ids
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	addr := os.Getenv("BRAINBOOST_COCKROACH_DB")
	db, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer db.Close()

	docs := getDocs(ctx, db)
	log.Printf("Found %d documents", len(docs))

	pine, err := pinecone.New()
	if err != nil {
		log.Fatal(err)
	}

	qdrant, err := vectordb_qdrant2.New()
	if err != nil {
		log.Fatal(err)
	}

	for inx, docId := range docs {
		chunks := getChunks(ctx, db, docId)
		log.Printf("[%-03d] Document %s has %d chunks", inx, docId, len(chunks))

		export, err := pine.Export(chunks)
		if err != nil {
			log.Fatal(err)
		}

		for _, chunk := range export {
			err = qdrant.Upsert([]*vectordb_qdrant2.Vector{{
				Id:           chunk.Id,
				DocumentId:   chunk.DocumentId,
				UserId:       chunk.UserId,
				CollectionId: chunk.CollectionId,
				Filename:     chunk.Filename,
				Text:         chunk.Text,
				Page:         chunk.Page,
				Vector:       chunk.Vector,
			}})
			if err != nil {
				log.Fatal(err)
			}
		}

		if inx > 600 {
			break
		}
	}
}
