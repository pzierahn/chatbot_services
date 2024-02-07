package chat

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"log"
)

func (service *Service) StartThread(ctx context.Context, prompt *pb.Prompt) (*pb.Thread, error) {
	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	if prompt.ModelOptions == nil {
		return nil, fmt.Errorf("options missing")
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

	model, err := service.getModel(prompt.ModelOptions.Model)
	if err != nil {
		return nil, err
	}

	resp, err := model.GenerateCompletion(ctx, &llm.GenerateRequest{
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

	completion := &pb.Thread{
		Id:           uuid.NewString(),
		CollectionId: prompt.CollectionId,
		Prompt:       prompt.Prompt,
		Completion:   resp.Text,
		References:   chunkData.ids,
		Scores:       chunkData.scores,
	}

	_ = service.storeThread(ctx, chatMessage{
		id:           completion.Id,
		userId:       userId,
		collectionId: prompt.CollectionId,
		prompt:       prompt.Prompt,
		completion:   completion.Completion,
		references:   chunkData.ids,
	})

	return completion, nil
}
