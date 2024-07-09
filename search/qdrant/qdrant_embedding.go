package qdrant

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/search"
	"log"
	"time"
)

const (
	batchSize        = 10
	parallelRequests = 10
)

var slots = make(chan struct{}, parallelRequests)

func init() {
	for inx := 0; inx < parallelRequests; inx++ {
		slots <- struct{}{}
	}
}

type embedding struct {
	id        []string
	embedding [][]float32
	tokens    uint32
	error     error
}

func (db *DB) createEmbeddings(ctx context.Context, fragments []*search.Fragment) (map[string][]float32, error) {
	results := make(chan *embedding, len(fragments))
	defer close(results)

	ctx, cnl := context.WithCancel(ctx)
	defer cnl()

	// Process 10 fragments in one go
	for start := 0; start < len(fragments); start += batchSize {
		end := min(start+batchSize, len(fragments))
		batch := fragments[start:end]

		// Start a goroutine for each document in parallel
		go func(batch []*search.Fragment) {
			select {
			case <-ctx.Done():
				// Abort if the context is canceled
				results <- &embedding{}
				return
			case <-slots:
				// Ensure the slot is released after the function returns
				defer func() { slots <- struct{}{} }()

				var inputs, ids []string
				for _, fragment := range batch {
					inputs = append(inputs, fragment.Text)
					ids = append(ids, fragment.Id)
				}

				// Allow up to 3 attempts to create an embedding
				for attempt := 2; attempt >= 0; attempt-- {
					result, err := db.embedding.CreateEmbedding(ctx, &llm.EmbeddingRequest{
						Inputs: inputs,
					})
					if err == nil {
						// Successfully created an embedding
						results <- &embedding{
							id:        ids,
							embedding: result.Embeddings,
							tokens:    result.Tokens,
						}
						break
					} else {
						// Failed to create an embedding. This can if too many requests are made in a short time.
						if attempt <= 0 {
							// Failed to create an embedding
							results <- &embedding{
								id:    ids,
								error: err,
							}
							break
						}

						// Wait for a short time before retrying
						time.Sleep(30 * time.Second)
					}
				}
			}
		}(batch)
	}

	received := 0
	embeddings := make(map[string][]float32)
	var err error

	for result := range results {
		received += len(result.id)
		// log.Printf("createEmbeddings: %d/%d", received, len(fragments))

		if err != nil {
			// Skip the remaining results if an error occurred
		} else if result.error != nil {
			// Record the error and cancel all other requests
			err = result.error
			cnl()

			log.Printf("error creating embeddings: %v", err)
		} else {
			// Record the embedding
			for inx, id := range result.id {
				embeddings[id] = result.embedding[inx]
			}
		}

		if received == len(fragments) {
			// All results have been received
			break
		}
	}

	if err != nil {
		return nil, err
	}

	return embeddings, nil
}
