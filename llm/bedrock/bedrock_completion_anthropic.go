package bedrock

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

type ClaudeMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ClaudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version,omitempty"`
	System           string          `json:"system,omitempty"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	Temperature      float64         `json:"temperature,omitempty"`
	TopP             float64         `json:"top_p,omitempty"`
	TopK             int             `json:"top_k,omitempty"`
	Messages         []ClaudeMessage `json:"messages,omitempty"`
}

type ResponseMessage struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

type ClaudeResponse struct {
	Id         string            `json:"id,omitempty"`
	Model      string            `json:"model,omitempty"`
	Content    []ResponseMessage `json:"content,omitempty"`
	Role       string            `json:"role,omitempty"`
	StopReason string            `json:"stop_reason,omitempty"`
	Type       string            `json:"type,omitempty"`
	Usage      ClaudeUsage       `json:"usage,omitempty"`
}

func (client *Client) generateCompletionAnthropic(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {

	var messages []ClaudeMessage

	for _, msg := range req.Messages {
		var role string
		switch msg.Type {
		case llm.MessageTypeUser:
			role = "user"
		case llm.MessageTypeBot:
			role = "assistant"
		}

		if len(messages) > 0 && messages[len(messages)-1].Role == role {
			messages[len(messages)-1].Content += "\n" + msg.Text
		} else {
			messages = append(messages, ClaudeMessage{
				Role:    role,
				Content: msg.Text,
			})
		}
	}

	request := ClaudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages:         messages,
		System:           req.SystemPrompt,
		MaxTokens:        1024,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	result, err := client.bedrock.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(req.Model),
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

	usage := llm.ModelUsage{
		UserId:       req.UserId,
		Model:        response.Model,
		InputTokens:  response.Usage.InputTokens,
		OutputTokens: response.Usage.OutputTokens,
	}

	client.usage.Track(ctx, usage)

	return &llm.GenerateResponse{
		Text:  strings.TrimSpace(response.Content[0].Text),
		Usage: usage,
	}, nil
}
