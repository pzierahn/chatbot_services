package mistral

import (
	"context"
	"github.com/gage-technologies/mistral-go"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

func (client *Client) GenerateCompletion(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	var messages []mistral.ChatMessage

	for _, msg := range req.Messages {
		switch msg.Type {
		case llm.MessageTypeSystem:
			messages = append(messages, mistral.ChatMessage{
				Content: msg.Text,
				Role:    mistral.RoleSystem,
			})
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

	client.usage.Track(ctx, llm.ModelUsage{
		UserId:           req.UserId,
		Model:            resp.Model,
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
	})

	return &llm.GenerateResponse{
		Text: response,
	}, nil
}
