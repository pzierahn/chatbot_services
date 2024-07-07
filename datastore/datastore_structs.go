package datastore

import (
	"github.com/google/uuid"
	"time"
)

// Function defines the function name and arguments
type Function struct {
	// Name of the function to call
	Name string `bson:"name,omitempty"`

	// Arguments to pass to the function
	Arguments map[string]any `bson:"arguments,omitempty"`
}

// ToolCall defines which tool to call
type ToolCall struct {
	// ID of the tool call
	CallID string `bson:"tool_call_id,omitempty"`

	// Define function to call
	Function Function `bson:"function,omitempty"`
}

// ToolResponse defines the response from the tool
type ToolResponse struct {
	// Calling tool ID
	CallID string `bson:"tool_call_id,omitempty"`

	// Tool response
	Content map[string]any `bson:"content,omitempty"`
}

// Message defines a message in the thread
type Message struct {
	// ID of the message
	Id uuid.UUID `bson:"_id,omitempty"`

	// Thread ID
	ThreadId uuid.UUID `bson:"thread_id,omitempty"`

	// User ID
	UserId string `bson:"user_id,omitempty"`

	// Role of the message
	Role string `bson:"role,omitempty"`

	// Timestamp of the message
	Timestamp time.Time `bson:"timestamp,omitempty"`

	// User or assistant message
	Content string `bson:"content,omitempty"`

	// Tool calls by assistant
	ToolCalls []ToolCall `bson:"tool_calls,omitempty"`

	// Tool calls response by tool
	ToolResponses []ToolResponse `bson:"tool_responses,omitempty"`
}
