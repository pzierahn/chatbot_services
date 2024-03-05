package openai

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
	"strings"
)

func (client *Client) GenerateCompletion(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	var messages []openai.ChatCompletionMessage

	if req.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemPrompt,
		})
	}

	for _, msg := range req.Messages {
		var role string
		switch msg.Type {
		case llm.MessageTypeUser:
			role = openai.ChatMessageRoleUser
		case llm.MessageTypeBot:
			role = openai.ChatMessageRoleAssistant
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Text,
		})
	}

	model, _ := strings.CutPrefix(req.Model, modelPrefix)

	resp, err := client.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       model,
			Temperature: req.Temperature,
			MaxTokens:   req.MaxTokens,
			Messages:    messages,
			N:           1,
			User:        req.UserId,
		},
	)
	if err != nil {
		return nil, err
	}

	usage := llm.ModelUsage{
		UserId:       req.UserId,
		Model:        resp.Model,
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
	}

	client.usage.Track(ctx, usage)

	return &llm.GenerateResponse{
		Text:  resp.Choices[0].Message.Content,
		Usage: usage,
	}, nil
}
