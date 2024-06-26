package chat

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"log"
)

func getDocumentText(doc *pb.Document) string {
	var text string

	for _, chunk := range doc.Chunks {
		text += chunk.Text
	}

	return text
}

func (service *Service) Completion(ctx context.Context, prompt *pb.CompletionRequest) (*pb.CompletionResponse, error) {
	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	if prompt.ModelOptions == nil {
		return nil, fmt.Errorf("options missing")
	}

	doc, err := service.docs.Get(ctx, &pb.DocumentID{Id: prompt.DocumentId})
	if err != nil {
		return nil, err
	}

	model, err := service.getModel(prompt.ModelOptions.Model)
	if err != nil {
		return nil, err
	}

	messages := []*llm.Message{{
		Type: llm.MessageTypeUser,
		Text: getDocumentText(doc) + "\n\n\n" + prompt.Prompt,
	}}

	resp, err := model.GenerateCompletion(ctx, &llm.GenerateRequest{
		SystemPrompt: "Be concise and short. Do not repeat parts of the prompt. Don't write any prefaces or introductions.",
		Messages:     messages,
		Model:        prompt.ModelOptions.Model,
		MaxTokens:    int(prompt.ModelOptions.MaxTokens),
		Temperature:  prompt.ModelOptions.Temperature,
		TopP:         prompt.ModelOptions.TopP,
		UserId:       userId,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	return &pb.CompletionResponse{
		Completion: resp.Text,
	}, nil
}
