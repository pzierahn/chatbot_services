package documents

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/vectordb"
	"log"
	"strings"
	"time"
)

const (
	parallelAgents = 50
)

type embedding struct {
	chunk     *pb.Chunk
	embedding []float32
	tokens    uint32
	error     error
}

type Embeddings map[string][]float32

func (service *Service) generateEmbeddings(ctx context.Context, userId string, doc *document) (Embeddings, error) {

	embeddings := make(Embeddings)

	results := make(chan *embedding, 1)
	defer close(results)

	queue := make(chan *pb.Chunk, len(doc.chunks))
	for _, chunk := range doc.chunks {
		if chunk.Id == "" {
			return nil, fmt.Errorf("chunk id is empty")
		}

		queue <- chunk
	}
	defer close(queue)

	for agent := 0; agent < parallelAgents; agent++ {
		go func(agent int) {
			for {
				select {
				case chunk, ok := <-queue:
					if !ok {
						return
					}

					text := strings.TrimSpace(chunk.Text)

					// Skip empty chunks
					if len(text) <= 0 {
						results <- &embedding{
							chunk: chunk,
							error: nil,
						}
						continue
					}

					ctx, cnl := context.WithTimeout(ctx, time.Second*5)
					resp, err := service.embeddings.CreateEmbeddings(ctx, &llm.EmbeddingRequest{
						Input:  text,
						UserId: userId,
						// Skip tracking for here to track it later
						SkipTracking: true,
					})
					cnl()

					if err != nil {
						// Some chunks may fail due to various reasons --> Retry chunk on error
						results <- &embedding{
							chunk: chunk,
							error: err,
						}

						break
					} else {
						results <- &embedding{
							chunk:     chunk,
							embedding: resp.Data,
							tokens:    uint32(resp.Tokens),
							error:     nil,
						}
					}

				case <-ctx.Done():
					return
				}
			}
		}(agent)
	}

	var tokenCount int
	var errorCount int
	var processed int

	for result := range results {
		if result.error != nil {
			if errorCount > 100 {
				return nil, result.error
			} else {
				queue <- result.chunk
				errorCount++
				continue
			}
		}

		if result.embedding != nil {
			embeddings[result.chunk.Id] = result.embedding
			tokenCount += int(result.tokens)
		}

		processed++
		if processed >= len(doc.chunks) {
			break
		}
	}

	// Pool embeddings tracking
	_, err := service.db.Exec(
		ctx,
		`INSERT INTO model_usages (user_id, model, input_tokens, output_tokens) 
			VALUES ($1, $2, $3, $4)`,
		userId,
		service.embeddings.GetEmbeddingModelName(),
		tokenCount,
		0,
	)
	if err != nil {
		log.Printf("Error tracking usage: %v", err)
	}

	return embeddings, nil
}

func (service *Service) insertEmbeddings(doc *document, embeddings Embeddings) error {
	var vectors []*vectordb.Vector

	for _, fragment := range doc.chunks {
		if fragment.Id == "" {
			return fmt.Errorf("fragment id is empty")
		}

		vector, ok := embeddings[fragment.Id]
		if !ok {
			return fmt.Errorf("missing embedding for fragment %s", fragment.Id)
		}

		vectors = append(vectors, &vectordb.Vector{
			Id:           fragment.Id,
			DocumentId:   doc.document.Id,
			CollectionId: doc.document.CollectionId,
			UserId:       doc.userId,
			Text:         fragment.Text,
			Vector:       vector,
		})
	}

	return service.vectorDB.Upsert(vectors)
}
