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

type FunctionCall func(ctx context.Context, input map[string]interface{}) (string, error)

// ToolDefinition defines a function that can be called by the assistant
type ToolDefinition struct {
	// Function name
	Name string

	// Description of the tool
	Description string

	// Parameters of the function
	Parameters ToolParameters

	// Call is the function to call
	Call FunctionCall
}

// ToolCall defines which tool to call
type ToolCall struct {
	// ID of the tool call
	CallID string `json:"id,omitempty" bson:"tool_call_id,omitempty"`

	// Name of the function to call
	Name string `json:"name,omitempty" bson:"name,omitempty"`

	// Arguments to pass to the function
	Arguments string `json:"arguments,omitempty" bson:"arguments,omitempty"`
}

// ToolResponse defines the response from the tool
type ToolResponse struct {
	// Calling tool ID
	CallID string `json:"id,omitempty" bson:"tool_call_id,omitempty"`

	// Tool response
	Content string `json:"content,omitempty" bson:"content,omitempty"`
}

// Message defines a message in the thread
type Message struct {
	// Role of the message
	Role string `json:"role,omitempty" bson:"role,omitempty"`

	// User or assistant message
	Content string `json:"content,omitempty" bson:"content,omitempty"`

	// Tool calls by assistant
	ToolCalls []ToolCall `json:"tool_calls,omitempty" bson:"tool_calls,omitempty"`

	// Tool calls response by tool
	ToolResponses []ToolResponse `json:"tool_responses,omitempty" bson:"tool_responses,omitempty"`
}

const (
	ToolUseAuto = "auto" // ToolUseAuto allows the bot to decide whether to call any provided tools or not
	ToolUseAny  = "any"  // ToolUseAny forces the bot to call any provided tools
	ToolUseTool = "tool" // ToolUseTool forces the bot to call a specific tool
	ToolUseNone = "none" // ToolUseNone forces the bot to not call any tools
)

// ToolChoice defines if the model should use a tool
type ToolChoice struct {
	Type string `json:"type,omitempty" bson:"type,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

// CompletionRequest defines the request to the completion API
type CompletionRequest struct {
	// SystemPrompt is the prompt for the system
	SystemPrompt string `json:"system_prompt,omitempty" bson:"system_prompt,omitempty"`

	// Messages in the thread
	Messages []*Message `json:"messages,omitempty" bson:"messages,omitempty"`

	// Model to use for completion
	Model string `json:"model,omitempty" bson:"model,omitempty"`

	// MaxTokens is the maximum number of tokens to generate
	MaxTokens int `json:"max_tokens,omitempty" bson:"max_tokens,omitempty"`

	// TopP is the nucleus sampling probability
	TopP float32 `json:"top_p,omitempty" bson:"top_p,omitempty"`

	// Temperature is the sampling temperature
	Temperature float32 `json:"temperature,omitempty" bson:"temperature,omitempty"`

	// UserId to prevent abuse
	UserId string `json:"user_id,omitempty" bson:"user_id,omitempty"`

	// ToolChoice defines if the user wants to use a tool
	ToolChoice *ToolChoice `json:"tool_choice,omitempty" bson:"tool_choice,omitempty"`

	// Tools to use for completion
	Tools []*ToolDefinition `json:"tools,omitempty" bson:"tools,omitempty"`
}

// CompletionResponse defines the response from the completion API
type CompletionResponse struct {
	// Message to return
	Messages []*Message `json:"messages,omitempty" bson:"messages,omitempty"`

	// Usage of the model
	Usage ModelUsage `json:"usage,omitempty" bson:"usage,omitempty"`
}

type Chat interface {
	Completion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
	ProvidesModel(model string) bool
}
