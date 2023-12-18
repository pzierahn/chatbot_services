package llm

import "context"

type GenerateRequest struct {
	Prompt      string
	Documents   []string
	Model       string
	MaxTokens   int
	TopP        float32
	Temperature float32
	UserId      string
}

type GenerateResponse struct {
	Text         string
	InputTokens  int
	OutputTokens int
}

type EmbeddingRequest struct {
	Input  string
	UserId string
}

type EmbeddingResponse struct {
	Data   []float32
	Tokens int
}

type LLM interface {
	CreateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
}
