package index

import (
	"context"
	"errors"
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
	"sync"
)

func (index Index) GetPagesWithEmbeddings(ctx context.Context, pages []string) ([]*database.PageEmbedding, error) {
	var mu sync.Mutex
	var embeddings []*database.PageEmbedding
	var errs []error

	var wg sync.WaitGroup
	wg.Add(len(pages))

	for inx, page := range pages {
		go func(inx int, page string) {
			defer wg.Done()

			log.Printf("Processing page %v", inx)
			page = strings.TrimSpace(page)
			if len(page) == 0 {
				return
			}

			resp, err := index.GPT.CreateEmbeddings(
				ctx,
				openai.EmbeddingRequestStrings{
					Model: openai.AdaEmbeddingV2,
					Input: []string{page},
				},
			)

			log.Printf("--> Processed page %v", inx)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs = append(errs, err)
				return
			}

			embeddings = append(embeddings, &database.PageEmbedding{
				Page:      inx,
				Text:      page,
				Embedding: resp.Data[0].Embedding,
			})
		}(inx, page)
	}

	wg.Wait()

	return embeddings, errors.Join(errs...)
}
