package llm

import (
	"context"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// ParametersProperties defines the properties of the parameters
type ParametersProperties struct {
	// Type of the parameter
	Type string

	// Description of the parameter
	Description string
}

// ToolParameters defines the input parameters for the function
type ToolParameters struct {
	// Type of the parameters
	Type string

	// Properties of the parameters
	Properties map[string]ParametersProperties

	// Required parameters
	Required []string
}

// ToolDefinition defines a function that can be called by the assistant
type ToolDefinition struct {
	// Function name
	Name string

	// Description of the tool
	Description string

	// Parameters of the function
	Parameters ToolParameters

	// Call is the function to call
	Call func(ctx context.Context, input map[string]interface{}) (string, error)
}

// Function defines the function name and arguments
type Function struct {
	// Name of the function to call
	Name string `json:"name,omitempty"`

	// Arguments to pass to the function
	Arguments string `json:"arguments,omitempty"`
}

// ToolCall defines which tool to call
type ToolCall struct {
	// ID of the tool call
	CallID string `json:"tool_call_id,omitempty"`

	// Define function to call
	Function Function `json:"function,omitempty"`
}

// ToolResponse defines the response from the tool
type ToolResponse struct {
	// Calling tool ID
	CallID string `json:"tool_call_id,omitempty"`

	// Tool response
	Content string `json:"content,omitempty"`
}

// Message defines a message in the thread
type Message struct {
	// Role of the message
	Role string `json:"role,omitempty"`

	// User or assistant message
	Content string `json:"content,omitempty"`

	// Tool calls by assistant
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`

	// Tool calls response by tool
	ToolResponses []ToolResponse `json:"tool_responses,omitempty"`
}

// CompletionRequest defines the request to the completion API
type CompletionRequest struct {
	// SystemPrompt is the prompt for the system
	SystemPrompt string `json:"system_prompt,omitempty"`

	// Messages in the thread
	Messages []*Message `json:"messages,omitempty"`

	// Model to use for completion
	Model string `json:"model,omitempty"`

	// MaxTokens is the maximum number of tokens to generate
	MaxTokens int `json:"max_tokens,omitempty"`

	// TopP is the nucleus sampling probability
	TopP float32 `json:"top_p,omitempty"`

	// Temperature is the sampling temperature
	Temperature float32 `json:"temperature,omitempty"`

	// UserId to prevent abuse
	UserId string `json:"user_id,omitempty"`
}

// CompletionResponse defines the response from the completion API
type CompletionResponse struct {
	// Message to return
	Message *Message `json:"message,omitempty"`

	// Usage of the model
	Usage ModelUsage `json:"usage,omitempty"`
}

type Chat interface {
	Completion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
	SetTools(tools []ToolDefinition)
	ProvidesModel(model string) bool
}
