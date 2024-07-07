package llm

import (
	"context"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"
)

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

type Function struct {
	// Name of the function to call
	Name string `json:"name,omitempty"`

	// Arguments to pass to the function
	Arguments string `json:"arguments,omitempty"`
}

type ToolCall struct {
	// ID of the tool call
	CallID string `json:"tool_call_id,omitempty"`

	// Define function to call
	Function Function `json:"function,omitempty"`
}

type ToolResponses struct {
	// Calling tool ID
	CallID string `json:"tool_call_id,omitempty"`

	// Tool response
	Content string `json:"content,omitempty"`
}

type Message struct {
	// Role of the message
	Role string `json:"role,omitempty"`

	// User or assistant message
	Content string `json:"content,omitempty"`

	// Tool calls by assistant
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`

	// Tool calls response by tool
	ToolResponses []ToolResponses `json:"tool_responses,omitempty"`
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
