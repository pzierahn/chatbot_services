package braingain

import (
	"context"
	"github.com/pzierahn/braingain/database"
	"github.com/sashabaranov/go-openai"
)

type Completion struct {
	Completion string
	Costs      Costs
}

type CompletionAugmentation struct {
	Completion
	Documents []database.ScorePoints
}

type Prompt struct {
	Prompt      string
	Model       string
	Background  []string
	Temperature float32
	MaxTokens   int
}

func (chat Chat) Chat(ctx context.Context, message Prompt) (*Completion, error) {

	var messages []openai.ChatCompletionMessage
	for _, text := range message.Background {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message.Prompt,
	})

	resp, err := chat.gpt.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       message.Model,
			Temperature: message.Temperature,
			MaxTokens:   message.MaxTokens,
			Messages:    messages,
			N:           1,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Completion{
		Completion: resp.Choices[0].Message.Content,
		Costs:      chat.calculateCosts(message.Model, resp.Usage),
	}, nil
}
