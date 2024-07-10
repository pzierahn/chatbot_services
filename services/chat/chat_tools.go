package chat

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/search"
	"log"
	"sort"
	"time"
)

type retrievalParameters struct {
	prompt        string
	userId        string
	collectionId  string
	fragmentCount uint32
	threshold     float32
}

type Sources struct {
	Items []*search.Result `json:"sources"`
}

type Document struct {
	Text string `json:"text"`
}

const (
	toolGetSources     = "get_sources"
	toolAttachDocument = "attach_document"
)

func (service *Service) getSourceTools(params retrievalParameters) *llm.ToolDefinition {
	return &llm.ToolDefinition{
		Name:        toolGetSources,
		Description: "Retrieve information from from the knowledge base. ",
		Parameters: llm.ToolParameters{
			Type: "object",
			Properties: map[string]llm.ParametersProperties{
				"prompt": {
					Type:        "string",
					Description: "A query or statement for which the information is requested.",
				},
			},
			Required: []string{
				"prompt",
			},
		},
		Call: func(ctx context.Context, parameters map[string]interface{}) (string, error) {
			query, ok := parameters["prompt"].(string)
			if !ok {
				return "", errors.New("prompt missing")
			}

			log.Printf("get_sources: %v", query)

			response, err := service.Search.Search(ctx, search.Query{
				UserId:       params.userId,
				CollectionId: params.collectionId,
				Query:        query,
				Limit:        params.fragmentCount,
				Threshold:    params.threshold,
			})
			if err != nil {
				return "", err
			}

			_ = service.Database.InsertModelUsage(ctx, &datastore.ModelUsage{
				Id:          uuid.New(),
				UserId:      params.userId,
				Timestamp:   time.Now(),
				ModelId:     response.Usage.ModelId,
				InputTokens: response.Usage.Tokens,
			})

			sources := response.Results

			// Group by document and sort by position
			sort.Slice(sources, func(i, j int) bool {
				if sources[i].DocumentId != sources[j].DocumentId {
					return sources[i].DocumentId < sources[j].DocumentId
				}
				return sources[i].Position < sources[j].Position
			})

			byt, err := json.Marshal(Sources{
				Items: sources,
			})
			if err != nil {
				return "", err
			}

			byt2, _ := json.MarshalIndent(sources, "", "  ")
			log.Printf("get_sources: %s", byt2)

			return string(byt), nil
		},
	}
}

func (service *Service) getAttachDocumentTool() *llm.ToolDefinition {
	return &llm.ToolDefinition{
		Name:        toolAttachDocument,
		Description: "Get the content of a document by its ID. Don't call this function!",
		Parameters: llm.ToolParameters{
			Type: "object",
			Properties: map[string]llm.ParametersProperties{
				"documentId": {
					Type:        "string",
					Description: "The ID of the document to retrieve.",
				},
			},
			Required: []string{
				"documentId",
			},
		},
		Call: func(ctx context.Context, parameters map[string]interface{}) (string, error) {
			return "", errors.New("don't call this function")
		},
	}
}
