package chat

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

func (service *Service) PostMessage(ctx context.Context, prompt *pb.Prompt) (*pb.Message, error) {
	log.Printf("PostMessage: %v", prompt)

	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
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

	threadId, err := uuid.Parse(prompt.ThreadId)
	if err != nil {
		return nil, err
	}

	messages, err := service.db.GetMessages(ctx, userId, threadId)
	if err != nil {
		return nil, err
	}

	//
	// Call the model
	//

	llmMessages, err := datastore.ToLLMMessages(messages)
	if err != nil {
		return nil, err
	}

	llmMessages = append(llmMessages, &llm.Message{
		Role:    "user",
		Content: prompt.Prompt,
	})

	datastorePrompt := &datastore.Message{
		Id:        uuid.New(),
		Role:      "user",
		Content:   prompt.Prompt,
		ThreadId:  threadId,
		UserId:    userId,
		Timestamp: time.Now(),
	}

	request := &llm.CompletionRequest{
		//SystemPrompt: "",
		Messages:    llmMessages,
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

	datastoreResponse, err := datastore.ToDatastoreMessage(response.Message)
	if err != nil {
		return nil, err
	}
	datastoreResponse.Id = uuid.New()
	datastoreResponse.ThreadId = threadId
	datastoreResponse.UserId = userId
	datastoreResponse.Timestamp = time.Now()

	err = service.db.AddMessages(ctx, []*datastore.Message{
		datastorePrompt,
		datastoreResponse,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Message{
		Id:         datastoreResponse.Id.String(),
		Prompt:     prompt.Prompt,
		Completion: response.Message.Content,
		Timestamp:  timestamppb.New(datastoreResponse.Timestamp),
	}, nil
}
