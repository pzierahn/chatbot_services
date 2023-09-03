package index

import (
	"context"
	"github.com/pzierahn/braingain/database"
	"github.com/pzierahn/braingain/pdf"
	"github.com/sashabaranov/go-openai"
	"log"
	"path/filepath"
	"strings"
)

func (index Index) File(ctx context.Context, filename string) {
	log.Printf("Filename: %s\n", filepath.Base(filename))

	pages, err := pdf.ReadPages(ctx, filename)
	if err != nil {
		log.Fatalf("could not read pages: %v", err)
	}

	log.Printf("--> Pages: %d\n", len(pages))

	// Create document
	docId, err := index.conn.CreateDocument(ctx, database.Document{
		Filename: filepath.Base(filename),
	})
	if err != nil {
		log.Fatalf("could not create document: %v", err)
	}

	for inx, page := range pages {
		log.Printf("--> %d/%d\n", inx+1, len(pages))

		page = strings.TrimSpace(page)
		if len(page) == 0 {
			continue
		}

		resp, err := index.ai.CreateEmbeddings(
			ctx,
			openai.EmbeddingRequestStrings{
				Model: openai.AdaEmbeddingV2,
				Input: []string{page},
			},
		)

		if err != nil {
			log.Fatalf("could not create embeddings: %v", err)
		}

		_, err = index.conn.Upsert(ctx, database.Point{
			Source:    docId,
			Page:      inx,
			Text:      page,
			Embedding: resp.Data[0].Embedding,
		})

		if err != nil {
			log.Fatalf("could not upsert points: %v", err)
		}
	}
}
