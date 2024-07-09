package llm

import "context"

type EmbeddingRequest struct {
	Inputs []string
	UserId string
}

type EmbeddingResponse struct {
	Embeddings [][]float32
	Tokens     uint32
	Model      string
}

type Embedding interface {
	CreateEmbedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
	GetEmbeddingDimension() int
	GetModelId() string
}
