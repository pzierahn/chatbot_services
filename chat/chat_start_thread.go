package chat

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

const (
	systemPromptQuote  = "When summarizing key points from the specified source, clearly indicate direct references by adding a reference right after the respective content. Only reference materials that are explicitly included in the context window."
	systemPromptNormal = "You are a helpful assistant. Answer in Markdown format without any code blocks. Give references to the source document if possible."
)

func (service *Service) StartThread(ctx context.Context, prompt *pb.ThreadPrompt) (*pb.Thread, error) {
	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	if prompt.ModelOptions == nil {
		return nil, fmt.Errorf("options missing")
	}

	var chunkData *chunks
	if len(prompt.DocumentIds) == 0 {
		chunkData, err = service.searchForContext(ctx, prompt)
	} else {
		chunkData, err = service.getDocumentsContext(ctx, userId, prompt)
	}
	if err != nil {
		return nil, err
	}

	messages := []*llm.Message{
		{
			Type: llm.MessageTypeSystem,
			Text: systemPromptQuote,
		},
	}

	for inx, doc := range chunkData.texts {
		messages = append(messages, &llm.Message{
			Type: llm.MessageTypeUser,
			Text: "Source: '" + chunkData.source[inx] + "'\n" + doc,
		})
	}

	// Add the prompt to the messages
	messages = append(messages, &llm.Message{
		Type: llm.MessageTypeUser,
		Text: prompt.Prompt,
	})

	model, err := service.getModel(prompt.ModelOptions.Model)
	if err != nil {
		return nil, err
	}

	resp, err := model.GenerateCompletion(ctx, &llm.GenerateRequest{
		Messages:    messages,
		Model:       prompt.ModelOptions.Model,
		MaxTokens:   int(prompt.ModelOptions.MaxTokens),
		Temperature: prompt.ModelOptions.Temperature,
		UserId:      userId,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	completion := &pb.Thread{
		Id:              uuid.NewString(),
		ReferenceIDs:    chunkData.ids,
		ReferenceScores: chunkData.scores,
		Messages: []*pb.Message{
			{
				Id:         uuid.NewString(),
				Prompt:     prompt.Prompt,
				Completion: resp.Text,
				Timestamp:  timestamppb.Now(),
			},
		},
	}

	err = service.storeThread(ctx, userId, prompt.CollectionId, completion)
	if err != nil {
		return nil, err
	}

	return completion, nil
}
