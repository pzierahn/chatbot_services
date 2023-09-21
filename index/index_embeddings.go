package index

import (
	"context"
	"errors"
	"github.com/pzierahn/brainboost/database"
	"github.com/sashabaranov/go-openai"
	"strings"
	"sync"
)

func (index Index) GetPagesWithEmbeddings(ctx context.Context, pages []string, ch ...chan<- Progress) ([]*database.PageEmbedding, error) {
	var mu sync.Mutex
	var embeddings []*database.PageEmbedding
	var errs []error

	var wg sync.WaitGroup
	wg.Add(len(pages))

	for inx, page := range pages {
		go func(inx int, page string) {
			defer wg.Done()
			defer func() {
				for _, c := range ch {
					c <- Progress{
						TotalPages:   len(pages),
						FinishedPage: inx,
					}
				}
			}()

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

	// TODO: Add model usage

	return embeddings, errors.Join(errs...)
}
