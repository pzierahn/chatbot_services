package voyageai

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
)

func (voyage *Client) CreateEmbedding(ctx context.Context, content *llm.EmbeddingRequest) (*llm.EmbeddingResponse, error) {
	request := &Request{
		Input:     content.Inputs,
		Model:     voyage.model,
		InputType: InputTypeDocument,
	}

	response, err := voyage.callAPI(ctx, request)
	if err != nil {
		return nil, err
	}

	embeddings := &llm.EmbeddingResponse{
		Model:  response.Model,
		Tokens: response.Usage.TotalTokens,
	}

	for _, output := range response.Data {
		embeddings.Embeddings = append(embeddings.Embeddings, output.Embedding)
	}

	return embeddings, nil
}

func (voyage *Client) GetEmbeddingDimension() int {
	switch voyage.model {
	case ModelVoyageLarge2:
		return DimensionVoyageLarge2
	case ModelVoyageLarge2Instruct:
		return DimensionVoyageLarge2Instruct
	default:
		return 0
	}
}
