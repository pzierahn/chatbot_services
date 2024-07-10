package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"github.com/pzierahn/chatbot_services/utils"
	"log"
	"strings"
	"time"
)

const (
	systemPromptLatex = "" +
		"You are a scientific research assistant. " +
		"Answer in Markdown format. " +
		"Quote sources with \\cite{document_id}."
	systemPromptNormal = "You are a helpful assistant. Answer in Markdown format."
)

// PostMessage is a gRPC endpoint that receives a prompt and returns a completion.
func (service *Service) PostMessage(ctx context.Context, prompt *pb.Prompt) (*pb.Message, error) {
	log.Printf("PostMessage: %v", utils.Prettify(prompt))

	userId, err := service.Auth.VerifyFunding(ctx)
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

	retrievalOptions := prompt.GetRetrievalOptions()
	if retrievalOptions == nil {
		return nil, fmt.Errorf("retrieval options missing")
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

		thread, err = service.Database.GetThread(ctx, userId, threadId)
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

	// Add the sources tool
	for _, docId := range prompt.Attachments {
		callId := uuid.New()

		documentId, err := uuid.Parse(docId)
		if err != nil {
			return nil, err
		}

		document, err := service.Database.GetDocument(ctx, userId, documentId)
		if err != nil {
			return nil, err
		}

		response, _ := json.Marshal(Document{
			Text: joinDocumentText(document),
		})

		messages = append(messages, []*llm.Message{
			{
				Role: llm.RoleAssistant,
				ToolCalls: []llm.ToolCall{{
					CallID: callId.String(),
					Function: llm.Function{
						Name:      toolAttachDocument,
						Arguments: fmt.Sprintf("{\"documentId\": \"%s\"}", documentId),
					},
				}},
			},
			{
				Role: llm.RoleUser,
				ToolResponses: []llm.ToolResponse{{
					CallID:  callId.String(),
					Content: string(response),
				}},
			},
		}...)
	}

	tools := []*llm.ToolDefinition{
		service.getAttachDocumentTool(),
		service.getSourceTools(retrievalParameters{
			prompt:        prompt.Prompt,
			userId:        userId,
			collectionId:  prompt.CollectionId,
			fragmentCount: retrievalOptions.Documents,
			threshold:     retrievalOptions.Threshold,
		}),
	}

	var systemPrompt string
	if len(prompt.Attachments) > 0 {
		// Attachment mode
		systemPrompt = systemPromptNormal
	} else {
		// Retrieval mode
		systemPrompt = systemPromptLatex
	}

	request := &llm.CompletionRequest{
		SystemPrompt: systemPrompt,
		Messages:     messages,
		Model:        modelOps.ModelId,
		MaxTokens:    int(modelOps.MaxTokens),
		TopP:         modelOps.TopP,
		Temperature:  modelOps.Temperature,
		UserId:       userId,
		Tools:        tools,
	}

	response, err := model.Completion(ctx, request)
	if err != nil {
		return nil, err
	}

	//
	// Save the response
	//

	thread.Messages = response.Messages
	err = service.Database.StoreThread(ctx, thread)
	if err != nil {
		return nil, err
	}

	_ = service.Database.InsertModelUsage(ctx, &datastore.ModelUsage{
		Id:           uuid.New(),
		UserId:       userId,
		Timestamp:    time.Now(),
		ModelId:      modelOps.ModelId,
		InputTokens:  response.Usage.InputTokens,
		OutputTokens: response.Usage.OutputTokens,
	})

	// Get the document names
	sources := getSources(response.Messages)

	for idx, source := range sources {
		docId, err := uuid.Parse(source.DocumentId)
		if err != nil {
			continue
		}

		docName, err := service.Database.GetDocumentName(ctx, userId, docId)
		if err != nil {
			continue
		}
		sources[idx].Name = docName
	}

	return &pb.Message{
		ThreadId:   thread.Id.String(),
		Prompt:     prompt.Prompt,
		Completion: thread.Messages[len(thread.Messages)-1].Content,
		Sources:    sources,
	}, nil
}

func joinDocumentText(document *datastore.Document) string {
	var parts []string

	for _, block := range document.Content {
		parts = append(parts, block.Text)
	}

	return strings.Join(parts, "\f")
}
