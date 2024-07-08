package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

// PostMessage is a gRPC endpoint that receives a prompt and returns a completion.
func (service *Service) PostMessage(ctx context.Context, prompt *pb.Prompt) (*pb.Message, error) {
	log.Printf("PostMessage: %v", prompt)

	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	//
	// Integrity check
	//

	collectionId, err := uuid.Parse(prompt.CollectionId)
	if err != nil {
		return nil, errors.New("invalid collection id")
	}

	modelOps := prompt.GetModelOptions()
	if modelOps == nil {
		return nil, fmt.Errorf("options missing")
	}

	model, err := service.getModel(modelOps.ModelId)
	if err != nil {
		return nil, err
	}

	//
	// Get the thread messages
	//

	var thread *datastore.Thread
	if prompt.ThreadId != "" {
		//
		// Get the thread from the database
		//

		threadId, err := uuid.Parse(prompt.ThreadId)
		if err != nil {
			return nil, err
		}

		thread, err = service.db.GetThread(ctx, userId, threadId)
		if err != nil {
			return nil, err
		}
	} else {
		//
		// Create a new thread
		//

		thread = &datastore.Thread{
			Id:           uuid.New(),
			UserId:       userId,
			CollectionId: collectionId,
			Timestamp:    time.Now(),
		}
	}

	//
	// Call the model
	//

	messages := append(thread.Messages, &llm.Message{
		Role:    llm.RoleUser,
		Content: prompt.Prompt,
	})

	request := &llm.CompletionRequest{
		//SystemPrompt: "",
		Messages:    messages,
		Model:       modelOps.ModelId,
		MaxTokens:   int(modelOps.MaxTokens),
		TopP:        modelOps.TopP,
		Temperature: modelOps.Temperature,
		UserId:      userId,
	}

	response, err := model.Completion(ctx, request)
	if err != nil {
		return nil, err
	}

	//
	// Save the response
	//

	thread.Messages = append(messages, response.Message)
	err = service.db.StoreThread(ctx, thread)
	if err != nil {
		return nil, err
	}

	return &pb.Message{
		ThreadId:   thread.Id.String(),
		Prompt:     prompt.Prompt,
		Completion: response.Message.Content,
		Timestamp:  timestamppb.Now(),
	}, nil
}
