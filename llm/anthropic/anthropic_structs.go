package anthropic

type ClaudeMessage struct {
	Role    string    `json:"role,omitempty"`
	Content []Content `json:"content,omitempty"`
}

type ToolChoice struct {
	Type string `json:"type,omitempty" bson:"type,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

type ClaudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version,omitempty"`
	System           string          `json:"system,omitempty"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	Temperature      float32         `json:"temperature,omitempty"`
	TopP             float32         `json:"top_p,omitempty"`
	TopK             int             `json:"top_k,omitempty"`
	Tools            []ClaudeTool    `json:"tools,omitempty"`
	ToolChoice       *ToolChoice     `json:"tool_choice,omitempty"`
	Messages         []ClaudeMessage `json:"messages,omitempty"`
}

type Content struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`

	// Function Parameters
	ID        string                 `json:"id,omitempty"`
	ToolUseId string                 `json:"tool_use_id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Input     map[string]interface{} `json:"input,omitempty"`
	Content   string                 `json:"content,omitempty"`
}

type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

type ClaudeResponse struct {
	Id         string      `json:"id,omitempty"`
	Model      string      `json:"model,omitempty"`
	Content    []Content   `json:"content,omitempty"`
	Role       string      `json:"role,omitempty"`
	StopReason string      `json:"stop_reason,omitempty"`
	Type       string      `json:"type,omitempty"`
	Usage      ClaudeUsage `json:"usage,omitempty"`
}

const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

const (
	ContentTypeText       = "text"
	ContentTypeToolUse    = "tool_use"
	ContentTypeToolResult = "tool_result"
)
