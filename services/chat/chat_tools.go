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

type documentParameters struct {
	userId string
}

type Sources struct {
	Items []*search.Result `json:"sources,omitempty" bson:"sources,omitempty"`
}

const (
	toolGetSources     = "get_sources"
	toolAttachDocument = "attach_document"
)

func (service *Service) getSourceTools(params retrievalParameters) *llm.ToolDefinition {
	return &llm.ToolDefinition{
		Name:        toolGetSources,
		Description: "Retrieves relevant knowledge and sources for a given query from a curated knowledge base. The search query should be clear, concise and containing key terms related to the desired information. The tool will return a list of document fragments with text and documentId. It should be used when the user asks about specific topics and requires precise information.",
		Parameters: llm.ToolParameters{
			Type: "object",
			Properties: map[string]llm.ParametersProperties{
				"query": {
					Type:        "string",
					Description: "A query or statement for which the information is requested.",
				},
			},
			Required: []string{
				"query",
			},
		},
		Call: func(ctx context.Context, parameters map[string]interface{}) (string, error) {
			query, ok := parameters["query"].(string)
			if !ok {
				return "", errors.New("query missing")
			}

			log.Printf("get_sources: \"%v\"", query)

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

			return string(byt), nil
		},
	}
}

func (service *Service) getDocumentById(ctx context.Context, userId, docId string) (string, error) {
	documentId, err := uuid.Parse(docId)
	if err != nil {
		return "", err
	}

	document, err := service.Database.GetDocument(ctx, userId, documentId)
	if err != nil {
		return "", err
	}

	sources := make([]*search.Result, len(document.Content))

	for idx, fragment := range document.Content {
		sources[idx] = &search.Result{
			Id:         fragment.Id.String(),
			Text:       fragment.Text,
			DocumentId: docId,
			Position:   fragment.Position,
		}
	}

	response, err := json.Marshal(Sources{
		Items: sources,
	})
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func (service *Service) getAttachDocumentTool(params documentParameters) *llm.ToolDefinition {
	return &llm.ToolDefinition{
		Name:        toolAttachDocument,
		Description: "Get the content of a document by its document_id. This tool should be used when the user asks for more information about a specific document. The tool will return the text of the document.",
		Parameters: llm.ToolParameters{
			Type: "object",
			Properties: map[string]llm.ParametersProperties{
				"document_id": {
					Type:        "string",
					Description: "The ID of the document to retrieve.",
				},
			},
			Required: []string{
				"document_id",
			},
		},
		Call: func(ctx context.Context, parameters map[string]interface{}) (string, error) {
			documentId, ok := parameters["document_id"].(string)
			if !ok {
				return "", errors.New("document_id missing")
			}

			log.Printf("attach_document: \"%v\"", documentId)

			return service.getDocumentById(ctx, params.userId, documentId)
		},
	}
}
