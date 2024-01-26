package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/openai"
	"github.com/pzierahn/chatbot_services/vectordb"
	"github.com/pzierahn/chatbot_services/vectordb/qdrant"
	"log"
	"os"
	"sync"
)

type Doc struct {
	Id           string
	UserId       string
	Filename     string
	CollectionId string
}

type Chunk struct {
	Id   string
	Text string
	Page int
}

// getDocs returns all documents from the database
func getDocs(ctx context.Context, db *pgxpool.Pool) []*Doc {
	rows, err := db.Query(ctx, `SELECT id, user_id, filename, collection_id  FROM documents`)
	if err != nil {
		log.Fatal(err)
	}

	var docs []*Doc
	for rows.Next() {
		var doc Doc
		err = rows.Scan(&doc.Id, &doc.UserId, &doc.Filename, &doc.CollectionId)
		if err != nil {
			log.Fatal(err)
		}
		docs = append(docs, &doc)
	}

	return docs
}

// getChunks returns all id, pages and texts for a given document
func getChunks(ctx context.Context, db *pgxpool.Pool, docId string) []*Chunk {
	rows, err := db.Query(ctx, `SELECT id, page, text FROM document_chunks WHERE document_id = $1`, docId)
	if err != nil {
		log.Fatal(err)
	}

	var chunks []*Chunk
	for rows.Next() {
		var chunk Chunk
		err = rows.Scan(&chunk.Id, &chunk.Page, &chunk.Text)
		if err != nil {
			log.Fatal(err)
		}
		chunks = append(chunks, &chunk)
	}

	return chunks
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

	gpt := openai.New()

	docs := getDocs(ctx, db)
	log.Printf("Found %d documents", len(docs))

	client, err := qdrant.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	//
	// Parallelize insert requests
	//

	var wg sync.WaitGroup
	ch := make(chan *vectordb.Vector, 1)

	for idx := 0; idx < 3; idx++ {
		go func() {
			for vc := range ch {
				embedding, err := gpt.CreateEmbeddings(ctx, &llm.EmbeddingRequest{
					Input:  vc.Text,
					UserId: vc.UserId,
				})
				if err != nil {
					log.Fatal(err)
				}

				vc.Vector = embedding.Data

				err = client.Upsert([]*vectordb.Vector{vc})
				if err != nil {
					log.Fatal(err)
				}

				wg.Done()
			}
		}()
	}

	//
	// Iterate over all documents and chunks
	//

	for inx, doc := range docs {
		if inx < 208 {
			continue
		}

		chunks := getChunks(ctx, db, doc.Id)
		wg.Add(len(chunks))

		percent := float64(inx) / float64(len(docs)) * 100
		log.Printf("[%d - %3.2f%%] Document %s has %d chunks", inx, percent, doc.Id, len(chunks))

		for _, chunk := range chunks {
			ch <- &vectordb.Vector{
				Id:           chunk.Id,
				DocumentId:   doc.Id,
				UserId:       doc.UserId,
				CollectionId: doc.CollectionId,
				Filename:     doc.Filename,
				Text:         chunk.Text,
				Page:         uint32(chunk.Page),
			}
		}

		//if inx > 100 {
		//	break
		//}
	}

	close(ch)

	log.Printf("Waiting for all requests to finish")
	wg.Wait()
}
