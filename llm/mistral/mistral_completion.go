package mistral

import (
	"context"
	"github.com/gage-technologies/mistral-go"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

func (client *Client) GenerateCompletion(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	var messages []mistral.ChatMessage

	if req.SystemPrompt != "" {
		messages = append(messages, mistral.ChatMessage{
			Content: req.SystemPrompt,
			Role:    mistral.RoleSystem,
		})
	}

	for _, msg := range req.Messages {
		switch msg.Type {
		case llm.MessageTypeUser:
			messages = append(messages, mistral.ChatMessage{
				Content: msg.Text,
				Role:    mistral.RoleUser,
			})
		case llm.MessageTypeBot:
			messages = append(messages, mistral.ChatMessage{
				Content: msg.Text,
				Role:    mistral.RoleAssistant,
			})
		}
	}

	model := strings.TrimPrefix(req.Model, prefix)
	resp, err := client.client.Chat(model, messages, &mistral.ChatRequestParams{
		Temperature: float64(req.Temperature),
		TopP:        float64(req.TopP),
		MaxTokens:   req.MaxTokens,
	})
	if err != nil {
		return nil, err
	}

	response := resp.Choices[0].Message.Content

	usage := llm.ModelUsage{
		UserId:       req.UserId,
		Model:        resp.Model,
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
	}

	client.usage.Track(ctx, usage)

	return &llm.GenerateResponse{
		Text:  response,
		Usage: usage,
	}, nil
}
