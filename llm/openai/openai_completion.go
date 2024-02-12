package openai

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
	"strings"
)

func (client *Client) GenerateCompletion(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant. " +
				"Answer in Markdown format without any code blocks." +
				"Don't include any latex code",
		},
	}

	for _, text := range req.Messages {
		var role string
		if text.Type == llm.MessageTypeBot {
			role = openai.ChatMessageRoleAssistant
		} else {
			role = openai.ChatMessageRoleUser
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: text.Text,
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

	client.trackUsage(ctx, llm.ModelUsage{
		UserId:           req.UserId,
		Model:            resp.Model,
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
	})

	return &llm.GenerateResponse{
		Text: resp.Choices[0].Message.Content,
	}, nil
}
