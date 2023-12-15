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

	if prompt.ModelOptions == nil {
		return nil, fmt.Errorf("options missing")
	}

	funding, err := service.account.HasFunding(ctx)
	if err != nil {
		return nil, err
	}

	if !funding {
		return nil, account.NoFundingError()
	}

	var chunkIds, fragments []string

	if len(prompt.Documents) > 0 {
		chunkIds, fragments, err = service.searchForContext(ctx, prompt)
	} else {
		chunkIds, fragments, err = service.getDocumentsContext(ctx, userId, prompt)
	}
	if err != nil {
		return nil, err
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Answer in Markdown format without any code blocks",
		},
	}

	for _, text := range fragments {
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
			Model:       prompt.ModelOptions.Model,
			Temperature: prompt.ModelOptions.Temperature,
			MaxTokens:   int(prompt.ModelOptions.MaxTokens),
			Messages:    messages,
			N:           1,
			User:        userId,
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
		Prompt:       prompt.Prompt,
		Text:         resp.Choices[0].Message.Content,
		References:   chunkIds,
	}

	_ = service.storeChatMessage(ctx, chatMessage{
		id:           completion.Id,
		userId:       userId,
		collectionId: prompt.CollectionId,
		prompt:       prompt.Prompt,
		completion:   completion.Text,
		references:   chunkIds,
	})

	return completion, nil
}
