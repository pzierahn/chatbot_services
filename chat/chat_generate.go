package chat

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/brainboost/account"
	"github.com/pzierahn/brainboost/llm"
	pb "github.com/pzierahn/brainboost/proto"
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

	var chunkData *chunks

	if len(prompt.Documents) == 0 {
		chunkData, err = service.searchForContext(ctx, prompt)
	} else {
		chunkData, err = service.getDocumentsContext(ctx, userId, prompt)
	}
	if err != nil {
		return nil, err
	}

	resp, err := service.completion.GenerateCompletion(ctx, &llm.GenerateRequest{
		Prompt:      prompt.Prompt,
		Documents:   chunkData.texts,
		Model:       prompt.ModelOptions.Model,
		MaxTokens:   int(prompt.ModelOptions.MaxTokens),
		Temperature: prompt.ModelOptions.Temperature,
		UserId:      userId,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	_, _ = service.account.CreateUsage(ctx, account.Usage{
		UserId: userId,
		Model:  prompt.ModelOptions.Model,
		Input:  uint32(resp.InputTokens),
		Output: uint32(resp.OutputTokens),
	})

	completion := &pb.ChatMessage{
		Id:           uuid.NewString(),
		CollectionId: prompt.CollectionId,
		Prompt:       prompt.Prompt,
		Text:         resp.Text,
		References:   chunkData.ids,
		Scores:       chunkData.scores,
	}

	_ = service.storeChatMessage(ctx, chatMessage{
		id:           completion.Id,
		userId:       userId,
		collectionId: prompt.CollectionId,
		prompt:       prompt.Prompt,
		completion:   completion.Text,
		references:   chunkData.ids,
	})

	return completion, nil
}
