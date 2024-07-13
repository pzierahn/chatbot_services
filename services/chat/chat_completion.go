package chat

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"log"
	"strings"
	"time"
)

// getDocumentText returns the text of a document as a single string.
func getDocumentText(document *datastore.Document) string {
	var parts []string

	for _, chunk := range document.Content {
		parts = append(parts, chunk.Text)
	}

	return strings.Join(parts, "\f")
}

// Completion retrieves a document from the database and sends it to the language model for completion with the given prompt.
func (service *Service) Completion(ctx context.Context, prompt *pb.CompletionRequest) (*pb.CompletionResponse, error) {
	userId, err := service.Auth.VerifyFunding(ctx)
	if err != nil {
		return nil, err
	}

	if prompt.ModelOptions == nil {
		return nil, fmt.Errorf("options missing")
	}

	docId, err := uuid.Parse(prompt.DocumentId)
	if err != nil {
		return nil, err
	}

	document, err := service.Database.GetDocument(ctx, userId, docId)
	if err != nil {
		return nil, err
	}

	model, err := service.getModel(prompt.ModelOptions.ModelId)
	if err != nil {
		return nil, err
	}

	messages := []*llm.Message{{
		Role:    llm.RoleUser,
		Content: getDocumentText(document) + "\n\n\n" + prompt.Prompt,
	}}

	response, err := model.Completion(ctx, &llm.CompletionRequest{
		SystemPrompt: "Be concise and short. Do not repeat parts of the prompt. Don't write any prefaces or introductions.",
		Messages:     messages,
		Model:        prompt.ModelOptions.ModelId,
		MaxTokens:    int(prompt.ModelOptions.MaxTokens),
		Temperature:  prompt.ModelOptions.Temperature,
		TopP:         prompt.ModelOptions.TopP,
		UserId:       userId,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	_ = service.Database.InsertModelUsage(ctx, &datastore.ModelUsage{
		Id:           uuid.New(),
		UserId:       userId,
		Timestamp:    time.Now(),
		ModelId:      response.Usage.Model,
		InputTokens:  response.Usage.InputTokens,
		OutputTokens: response.Usage.OutputTokens,
	})

	return &pb.CompletionResponse{
		Completion: response.Messages[1].Content,
	}, nil
}
