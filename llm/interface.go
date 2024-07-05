package llm

import (
	"context"
	"log"
)

const (
	MessageTypeUser = iota
	MessageTypeAssistant
)

type Message struct {
	Type int
	Text string
}

type GenerateRequest struct {
	SystemPrompt string
	Messages     []*Message
	Model        string
	MaxTokens    int
	TopP         float32
	Temperature  float32
	UserId       string
}

type GenerateResponse struct {
	Text  string
	Usage ModelUsage
}

type EmbeddingRequest struct {
	Input        string
	UserId       string
	SkipTracking bool
}

type EmbeddingResponse struct {
	Data   []float32
	Tokens int
}

type ModelUsage struct {
	Model        string `json:"model,omitempty"`
	UserId       string `json:"user_id,omitempty"`
	InputTokens  int    `json:"prompt_tokens,omitempty"`
	OutputTokens int    `json:"completion_tokens,omitempty"`
}

type PricePer1000Tokens struct {
	Input  float32
	Output float32
}

// Cost returns the cost in cents of a model usage
func (price *PricePer1000Tokens) Cost(input, output uint32) (cost uint32) {
	cost += uint32(float32(input) * price.Input)
	cost += uint32(float32(output) * price.Output)
	return cost / 10
}

type Usage interface {
	Track(ctx context.Context, usage ModelUsage)
}

type Embedding interface {
	CreateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
	GetEmbeddingModelName() string
}

type Completion interface {
	GenerateCompletion(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
	ProvidesModel(model string) bool
}

type DummyTracker struct {
	PrintUsage bool
}

func (dummy DummyTracker) Track(ctx context.Context, usage ModelUsage) {
	if dummy.PrintUsage {
		log.Printf("Usage: %v", usage)
	}
}
