package openai

import (
	"context"
	"github.com/pzierahn/brainboost/llm"
	"github.com/sashabaranov/go-openai"
)

func (client *Client) Generate(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Answer in Markdown format without any code blocks",
		},
	}

	for _, text := range req.Documents {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Prompt,
	})

	var model string
	if req.Model != "" {
		model = req.Model
	} else {
		model = openai.GPT4TurboPreview
	}

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

	return &llm.GenerateResponse{
		Text:         resp.Choices[0].Message.Content,
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
	}, nil
}
