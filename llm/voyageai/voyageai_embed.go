package voyageai

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
)

func (voayage *Client) CreateEmbedding(ctx context.Context, content *llm.EmbeddingRequest) (*llm.EmbeddingResponse, error) {
	request := &Request{
		Input:     content.Inputs,
		Model:     voayage.model,
		InputType: InputTypeDocument,
	}

	response, err := voayage.callAPI(ctx, request)
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
