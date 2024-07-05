package llm

import (
	"context"
	"time"
)

const (
	MessageTypeUser = iota
	MessageTypeAssistant
	MessageTypeTool
)

type Function struct {
	Name      string
	Arguments string
}

type ToolCall struct {
	Id       string
	Type     string
	Function Function
}

type ParametersProperties struct {
	Type        string
	Description string
}

type ToolParameters struct {
	Type       string
	Properties map[string]ParametersProperties
	Required   []string
}

type ToolDefinition struct {
	Name        string
	Description string
	Parameters  ToolParameters
	Call        func(ctx context.Context, input map[string]interface{}) (string, error)
}

type Message struct {
	Role      int        `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Timestamp time.Time  `json:"timestamp,omitempty"`
}

type CompletionRequest struct {
	SystemPrompt string     `json:"system_prompt,omitempty"`
	Messages     []*Message `json:"messages,omitempty"`
	Model        string     `json:"model,omitempty"`
	MaxTokens    int        `json:"max_tokens,omitempty"`
	TopP         float32    `json:"top_p,omitempty"`
	Temperature  float32    `json:"temperature,omitempty"`
	UserId       string     `json:"user_id,omitempty"`
}

type CompletionResponse struct {
	Message *Message   `json:"message,omitempty"`
	Usage   ModelUsage `json:"usage,omitempty"`
}

type Chat interface {
	Completion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
	SetTools(tools []ToolDefinition)
	ProvidesModel(model string) bool
}
