package chat

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

func (service *Service) PostMessage(ctx context.Context, prompt *pb.Prompt) (*pb.Message, error) {
	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	if prompt.ModelOptions == nil {
		return nil, fmt.Errorf("options missing")
	}

	references, err := service.getReferences(ctx, userId, prompt.ThreadID)
	if err != nil {
		return nil, err
	}

	messages := []*llm.Message{
		{
			Type: llm.MessageTypeSystem,
			Text: systemPromptQuote,
		},
	}

	for _, ref := range references {
		messages = append(messages, &llm.Message{
			Type: llm.MessageTypeUser,
			Text: ref.Text,
		})
	}

	chatMessages, err := service.getThreadMessages(ctx, userId, prompt.ThreadID)
	if err != nil {
		return nil, err
	}
	for _, msg := range chatMessages {
		messages = append(messages, []*llm.Message{{
			Type: llm.MessageTypeUser,
			Text: msg.Prompt,
		}, {
			Type: llm.MessageTypeBot,
			Text: msg.Completion,
		}}...)
	}

	messages = append(messages, &llm.Message{
		Type: llm.MessageTypeUser,
		Text: prompt.Prompt,
	})

	model, err := service.getModel(prompt.ModelOptions.Model)
	if err != nil {
		log.Printf("Error fetching model: %v", err)
		return nil, err
	}

	completion, err := model.GenerateCompletion(ctx, &llm.GenerateRequest{
		Messages:    messages,
		Model:       prompt.ModelOptions.Model,
		MaxTokens:   1024,
		Temperature: prompt.ModelOptions.Temperature,
		TopP:        prompt.ModelOptions.TopP,
		UserId:      userId,
	})
	if err != nil {
		log.Printf("Error generating completion: %v", err)
		return nil, err
	}

	message := &pb.Message{
		Prompt:     prompt.Prompt,
		Completion: completion.Text,
		Timestamp:  timestamppb.Now(),
	}

	err = service.db.QueryRow(
		ctx,
		`INSERT INTO thread_messages (user_id, thread_id, prompt, completion, created_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING (id)`,
		userId,
		prompt.ThreadID,
		prompt.Prompt,
		completion.Text,
		utils.ProtoToTime(message.Timestamp)).
		Scan(&message.Id)
	if err != nil {
		log.Printf("Error inserting message: %v", err)
		return nil, err
	}

	return message, err
}
