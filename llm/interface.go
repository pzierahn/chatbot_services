package llm

import "context"

const (
	MessageTypeUser = iota
	MessageTypeBot
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

type Usage interface {
	Track(ctx context.Context, usage ModelUsage)
}

type Embedding interface {
	CreateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
	GetModelName() string
}

type Completion interface {
	GenerateCompletion(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
	ProvidesModel(model string) bool
}
