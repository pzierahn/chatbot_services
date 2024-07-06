package llm

import (
	"context"
)

const (
	MessageTypeUser      = "user"
	MessageTypeAssistant = "assistant"
	MessageTypeTool      = "tool"
)

type Function struct {
	Name      string
	Arguments string
}

type ToolCall struct {
	Id       string
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
	Role       string     `json:"role,omitempty"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
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
