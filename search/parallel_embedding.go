package search

import (
	"context"
	"errors"
	"github.com/pzierahn/chatbot_services/llm"
	"log"
	"time"
)

// ParallelEmbedding defines an embedding engine that processes multiple fragments in parallel.
type ParallelEmbedding struct {
	engine    llm.Embedding // Engine defines the embedding engine
	batchSize int           // BatchSize defines how many fragments are processed in one request
	slots     chan struct{} // Slots defines how many requests can be made in parallel
}

// EmbeddingResponse defines the response of a parallel embedding request.
type EmbeddingResponse struct {
	Embeddings map[string][]float32
	Usage      Usage
}

type embedding struct {
	id        []string
	embedding [][]float32
	tokens    uint32
	error     error
}

// NewParallelEmbedding creates a new parallel embedding engine.
func NewParallelEmbedding(engine llm.Embedding, agents, batchSize int) *ParallelEmbedding {
	parallelEmbedding := &ParallelEmbedding{
		engine:    engine,
		batchSize: batchSize,
		slots:     make(chan struct{}, agents),
	}

	// Slots defines how many requests can be made in parallel
	for inx := 0; inx < agents; inx++ {
		parallelEmbedding.slots <- struct{}{}
	}

	return parallelEmbedding
}

// Close releases all resources.
func (engine *ParallelEmbedding) Close() {
	close(engine.slots)
}

// CreateEmbeddings creates embeddings for multiple fragments in parallel.
func (engine *ParallelEmbedding) CreateEmbeddings(ctx context.Context, fragments []*Fragment) (*EmbeddingResponse, error) {
	results := make(chan *embedding, len(fragments))
	defer close(results)

	ctx, cnl := context.WithCancel(ctx)
	defer cnl()

	// Process 10 fragments in one request
	for start := 0; start < len(fragments); start += engine.batchSize {
		end := min(start+engine.batchSize, len(fragments))
		batch := fragments[start:end]

		// Start a goroutine for each document in parallel
		go func(batch []*Fragment) {
			select {
			case <-ctx.Done():
				// Abort if the context is canceled
				results <- &embedding{}
				return
			case _, ok := <-engine.slots:
				if !ok {
					// Abort if the slot is not available
					results <- &embedding{
						error: errors.New("slot not available"),
					}
					return
				}

				// Ensure the slot is released after the function returns
				defer func() { engine.slots <- struct{}{} }()

				var inputs, ids []string
				for _, fragment := range batch {
					inputs = append(inputs, fragment.Text)
					ids = append(ids, fragment.Id)
				}

				// Allow up to 3 attempts to create an embedding
				for attempt := 2; attempt >= 0; attempt-- {
					result, err := engine.engine.CreateEmbedding(ctx, &llm.EmbeddingRequest{
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
	response := &EmbeddingResponse{
		Embeddings: make(map[string][]float32),
		Usage: Usage{
			ModelId: engine.engine.GetModelId(),
			Tokens:  0,
		},
	}
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
				response.Embeddings[id] = result.embedding[inx]
				response.Usage.Tokens += result.tokens
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

	return response, nil
}
