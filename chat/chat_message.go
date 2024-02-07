package chat

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (service *Service) PostMessage(ctx context.Context, prompt *pb.Prompt) (*pb.Message, error) {
	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	if prompt.ModelOptions == nil {
		return nil, fmt.Errorf("options missing")
	}

	model, err := service.getModel(prompt.ModelOptions.Model)
	if err != nil {
		return nil, err
	}

	completion, err := model.GenerateCompletion(ctx, &llm.GenerateRequest{
		Prompt:      prompt.Prompt,
		Documents:   nil,
		Model:       prompt.ModelOptions.Model,
		MaxTokens:   1024,
		Temperature: prompt.ModelOptions.Temperature,
		UserId:      userId,
	})
	if err != nil {
		return nil, err
	}

	message := &pb.Message{
		Prompt:       prompt.Prompt,
		Completion:   completion.Text,
		ModelOptions: prompt.ModelOptions,
	}

	var createdAt time.Time

	err = service.db.QueryRow(
		ctx,
		`INSERT INTO thread_messages (user_id, thread_id, prompt, completion)
			VALUES ($1, $2, $3, $4)
			RETURNING (id, created_at)`,
		userId, prompt.ThreadID, prompt.Prompt, completion.Text).
		Scan(&message.Id, &createdAt)
	if err != nil {
		return nil, err
	}

	message.Timestamp = timestamppb.New(createdAt)

	return message, err
}
