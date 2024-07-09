package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/vectordb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type Sources struct {
	Items []*vectordb.SearchResult `json:"sources"`
}

// PostMessage is a gRPC endpoint that receives a prompt and returns a completion.
func (service *Service) PostMessage(ctx context.Context, prompt *pb.Prompt) (*pb.Message, error) {
	log.Printf("PostMessage: %v", prompt)

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

	request := &llm.CompletionRequest{
		//SystemPrompt: "",
		Messages:    messages,
		Model:       modelOps.ModelId,
		MaxTokens:   int(modelOps.MaxTokens),
		TopP:        modelOps.TopP,
		Temperature: modelOps.Temperature,
		UserId:      userId,
		Tools: []llm.ToolDefinition{{
			Name:        "get_sources",
			Description: "Retrieves the sources for the prompt. The prompt should be optimized for embedding retrieval. The tool will return a list of sources in JSON format with the following fields: SourceID, Content.",
			Parameters: llm.ToolParameters{
				Type: "object",
				Properties: map[string]llm.ParametersProperties{
					"prompt": {
						Type:        "string",
						Description: "The topic for which to retrieve sources. The prompt should be optimized for embedding retrieval.",
					},
				},
				Required: []string{"prompt"},
			},
			Call: func(ctx context.Context, parameters map[string]interface{}) (string, error) {
				query, ok := parameters["prompt"].(string)
				if !ok {
					return "", errors.New("prompt missing")
				}

				log.Printf("get_sources: %v", query)

				search, err := service.Search.Search(ctx, vectordb.SearchQuery{
					UserId:       userId,
					CollectionId: prompt.CollectionId,
					Query:        query,
					Limit:        retrievalOptions.Documents,
					Threshold:    retrievalOptions.Threshold,
				})
				if err != nil {
					return "", err
				}

				byt, err := json.Marshal(Sources{
					Items: search,
				})
				if err != nil {
					return "", err
				}

				byt2, _ := json.MarshalIndent(search, "", "  ")
				log.Printf("get_sources: %s", byt2)

				return string(byt), nil
			},
		}},
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

	return &pb.Message{
		ThreadId:   thread.Id.String(),
		Prompt:     prompt.Prompt,
		Completion: thread.Messages[len(thread.Messages)-1].Content,
		Timestamp:  timestamppb.Now(),
		References: getSources(response.Messages),
	}, nil
}

func getSources(messages []*llm.Message) map[string]float32 {
	if len(messages) < 3 {
		return nil
	}

	sources := make(map[string]float32)

	for idx := len(messages) - 2; idx > 0; idx-- {
		message := messages[idx]
		//if message.Role == llm.RoleUser && len(message.ToolCalls) == 0 {
		//	break
		//}

		isSourceCall := make(map[string]bool)
		for _, toolCall := range messages[idx-1].ToolCalls {
			if toolCall.Function.Name == "get_sources" {
				isSourceCall[toolCall.CallID] = true
			}
		}

		for _, toolResponse := range message.ToolResponses {
			if isSourceCall[toolResponse.CallID] {
				var source Sources

				err := json.Unmarshal([]byte(toolResponse.Content), &source)
				if err != nil {
					continue
				}

				for _, item := range source.Items {
					sources[item.Id] = item.Score
				}
			}
		}
	}

	return sources
}
