package chat

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
)

func getDocumentText(doc *pb.Document) string {
	var text string

	for _, chunk := range doc.Chunks {
		text += chunk.Text
	}

	return text
}

func getDocumentTitle(doc *pb.Document) string {
	switch doc.Metadata.Data.(type) {
	case *pb.DocumentMetadata_Web:
		return doc.Metadata.GetWeb().Title
	case *pb.DocumentMetadata_File:
		return doc.Metadata.GetFile().Filename
	}

	return ""
}

func (service *Service) BatchChat(ctx context.Context, req *pb.BatchRequest) (*pb.BatchResponse, error) {
	userId, err := service.Verify(ctx)
	if err != nil {
		return nil, err
	}

	model, err := service.getModel(req.ModelOptions.Model)
	if err != nil {
		return nil, err
	}

	response := &pb.BatchResponse{
		DocumentIds: req.DocumentIds,
		Prompts:     req.Prompts,
		PromptTitle: make([]string, len(req.Prompts)),
	}

	for promptIdx, prompt := range req.Prompts {
		resp, err := model.GenerateCompletion(ctx, &llm.GenerateRequest{
			Messages: []*llm.Message{{
				Type: llm.MessageTypeUser,
				Text: fmt.Sprintf("Create a table column name from this prompt in snake case: '%s'. "+
					"Do it without any additional words or explanations.", prompt),
			}},
			Model:       req.ModelOptions.Model,
			MaxTokens:   10,
			TopP:        0,
			Temperature: 0,
			UserId:      userId,
		})
		if err != nil {
			return nil, err
		}

		response.PromptTitle[promptIdx] = resp.Text
	}

	for docIdx, docID := range req.DocumentIds {
		doc, err := service.docs.Get(ctx, &pb.DocumentID{Id: docID})
		if err != nil {
			return nil, err
		}

		text := getDocumentText(doc)
		title := getDocumentTitle(doc)

		for promptIdx, prompt := range req.Prompts {
			completion, err := model.GenerateCompletion(ctx, &llm.GenerateRequest{
				Model:       req.ModelOptions.Model,
				MaxTokens:   int(req.ModelOptions.MaxTokens),
				Temperature: req.ModelOptions.Temperature,
				TopP:        req.ModelOptions.TopP,
				UserId:      userId,
				Messages: []*llm.Message{{
					Text: text + "\n\n\n" + prompt,
				}},
			})
			if err != nil {
				return nil, err
			}

			response.Items = append(response.Items, &pb.BatchResponse_Completion{
				DocumentId:    uint32(docIdx),
				DocumentTitle: title,
				Prompt:        uint32(promptIdx),
				Completion:    completion.Text,
			})
		}
	}

	return response, nil
}
