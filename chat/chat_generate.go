package chat

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/account"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/sashabaranov/go-openai"
	"log"
)

func (service *Service) Chat(ctx context.Context, prompt *pb.Prompt) (*pb.ChatMessage, error) {
	userId, err := service.auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}

	if prompt.Options == nil {
		return nil, fmt.Errorf("options missing")
	}

	var bg *chatContext
	if prompt.Documents == nil || len(prompt.Documents) == 0 {
		bg, err = service.getSourceFromDB(ctx, prompt)
	} else {
		bg, err = service.getBackgroundFromPrompt(ctx, userId, prompt)
	}

	if err != nil {
		return nil, err
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Always answer in Markdown format",
		},
	}

	for _, text := range bg.fragments {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt.Prompt,
	})

	resp, err := service.gpt.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       prompt.Options.Model,
			Temperature: prompt.Options.Temperature,
			MaxTokens:   int(prompt.Options.MaxTokens),
			Messages:    messages,
			N:           1,
			User:        userId.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = service.account.CreateUsage(ctx, account.Usage{
		UserId: userId,
		Model:  resp.Model,
		Input:  uint32(resp.Usage.PromptTokens),
		Output: uint32(resp.Usage.CompletionTokens),
	})
	if err != nil {
		log.Printf("Chat: error %v", err)
	}

	completion := &pb.ChatMessage{
		Id:           uuid.NewString(),
		CollectionId: prompt.CollectionId,
		Prompt:       prompt,
		Text:         resp.Choices[0].Message.Content,
		Documents:    bg.docs,
	}

	_ = service.storeChatMessage(ctx, chatMessage{
		id:           completion.Id,
		userId:       userId,
		collectionId: prompt.CollectionId,
		prompt:       prompt.Prompt,
		completion:   completion.Text,
		references:   bg.pageIds,
	})

	return completion, nil
}
