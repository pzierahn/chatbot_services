package llm

import "context"

const (
	EmbeddingTypeQuery    = "query"
	EmbeddingTypeDocument = "document"
)

type EmbeddingRequest struct {
	Inputs []string
	UserId string
	Type   string
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
