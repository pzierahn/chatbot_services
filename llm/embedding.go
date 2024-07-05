package llm

import "context"

type EmbeddingRequest struct {
	Input        string
	UserId       string
	SkipTracking bool
}

type EmbeddingResponse struct {
	Data   []float32
	Tokens int
}

type Embedding interface {
	Create(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
	GetEmbeddingModelName() string
}
