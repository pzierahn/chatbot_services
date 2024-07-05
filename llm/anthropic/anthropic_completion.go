package anthropic

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
	"time"
)

type ClaudeMessage struct {
	Role    string    `json:"role,omitempty"`
	Content []Content `json:"content,omitempty"`
}

type ClaudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version,omitempty"`
	System           string          `json:"system,omitempty"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	Temperature      float32         `json:"temperature,omitempty"`
	TopP             float32         `json:"top_p,omitempty"`
	TopK             int             `json:"top_k,omitempty"`
	Tools            []ClaudeTool    `json:"tools,omitempty"`
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

func (client *Client) invokeRequest(model string, req *ClaudeRequest) (*ClaudeResponse, error) {
	body, _ := json.Marshal(req)
	result, err := client.bedrock.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(model),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return nil, err
	}

	var response ClaudeResponse
	err = json.Unmarshal(result.Body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	var messages []ClaudeMessage

	for _, msg := range req.Messages {
		var role string
		switch msg.Role {
		case llm.MessageTypeUser:
			role = ChatMessageRoleUser
		case llm.MessageTypeAssistant:
			role = ChatMessageRoleAssistant
		case llm.MessageTypeTool:
			role = ChatMessageRoleUser
		}

		messages = append(messages, ClaudeMessage{
			Role: role,
			Content: []Content{{
				Type: "text",
				Text: msg.Content,
			}},
		})
	}

	request := ClaudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages:         messages,
		System:           req.SystemPrompt,
		MaxTokens:        req.MaxTokens,
		TopP:             req.TopP,
		Temperature:      req.Temperature,
		Tools:            client.getTools(),
	}

	response, err := client.invokeRequest(req.Model, &request)
	if err != nil {
		return nil, err
	}

	usage := llm.ModelUsage{
		UserId:       req.UserId,
		Model:        response.Model,
		InputTokens:  response.Usage.InputTokens,
		OutputTokens: response.Usage.OutputTokens,
	}

	if response.StopReason == "tool_use" {
		messages = append(messages, ClaudeMessage{
			Role:    ChatMessageRoleAssistant,
			Content: response.Content,
		})

		for _, message := range response.Content {
			if message.Type == "tool_use" {
				result, err := client.callTool(ctx, message.Name, message.Input)
				if err != nil {
					return nil, err
				}

				messages = append(messages, ClaudeMessage{
					Role: ChatMessageRoleUser,
					Content: []Content{{
						Type:      "tool_result",
						ToolUseId: message.ID,
						Content:   result,
					}},
				})
			}
		}

		request.Messages = messages
		response, err = client.invokeRequest(req.Model, &request)
		if err != nil {
			return nil, err
		}

		usage.OutputTokens += response.Usage.OutputTokens
		usage.InputTokens += response.Usage.InputTokens
	}

	return &llm.CompletionResponse{
		Message: &llm.Message{
			Role:      llm.MessageTypeAssistant,
			Content:   strings.TrimSpace(response.Content[0].Text),
			Timestamp: time.Now(),
		},
		Usage: usage,
	}, nil
}
