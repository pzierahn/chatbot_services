package qdrant

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/vectordb"
	"log"
	"time"
)

const parallelRequests = 10

var slots = make(chan struct{}, parallelRequests)

func init() {
	for inx := 0; inx < parallelRequests; inx++ {
		slots <- struct{}{}
	}
}

type embedding struct {
	id        string
	embedding []float32
	tokens    uint32
	error     error
}

func (db *DB) createEmbeddings(ctx context.Context, fragments []*vectordb.Fragment) (map[string][]float32, error) {
	results := make(chan *embedding, len(fragments))
	defer close(results)

	ctx, cnl := context.WithCancel(ctx)
	defer cnl()

	for _, fragment := range fragments {
		// Start a goroutine for each document in parallel
		go func(fragment *vectordb.Fragment) {
			select {
			case <-ctx.Done():
				// Abort if the context is canceled
				results <- &embedding{}
				return
			case <-slots:
				// Ensure the slot is released after the function returns
				defer func() { slots <- struct{}{} }()

				// Allow up to 3 attempts to create an embedding
				for attempt := 2; attempt >= 0; attempt-- {
					result, err := db.embedding.CreateEmbedding(ctx, &llm.EmbeddingRequest{
						Inputs: []string{fragment.Text},
						UserId: "todo",
					})
					if err == nil {
						// Successfully created an embedding
						results <- &embedding{
							id:        fragment.Id,
							embedding: result.Embeddings[0],
							tokens:    result.Tokens,
						}
						break
					} else {
						// Failed to create an embedding. This can if too many requests are made in a short time.
						if attempt <= 0 {
							// Failed to create an embedding
							results <- &embedding{
								id:    fragment.Id,
								error: err,
							}
							break
						}

						// Wait for a short time before retrying
						time.Sleep(30 * time.Second)
					}
				}
			}
		}(fragment)
	}

	received := 0
	embeddings := make(map[string][]float32)
	var err error

	for result := range results {
		received++
		log.Printf("createEmbeddings: %d/%d", received, len(fragments))

		if err != nil {
			// Skip the remaining results if an error occurred
		} else if result.error != nil {
			// Record the error and cancel all other requests
			err = result.error
			log.Printf("createEmbeddings: %v", err)
			cnl()
		} else {
			// Record the embedding
			embeddings[result.id] = result.embedding
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
