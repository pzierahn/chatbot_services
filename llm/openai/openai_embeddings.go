package openai

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
)

func (client *Client) CreateEmbedding(ctx context.Context, req *llm.EmbeddingRequest) (*llm.EmbeddingResponse, error) {
	resp, err := client.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Model: client.embeddingModel,
			Input: req.Inputs,
			User:  req.UserId,
		},
	)
	if err != nil {
		return nil, err
	}

	results := &llm.EmbeddingResponse{
		Embeddings: make([][]float32, len(resp.Data)),
		Model:      client.GetModelId(),
		Tokens:     uint32(resp.Usage.PromptTokens),
	}

	for idx, item := range resp.Data {
		results.Embeddings[idx] = item.Embedding
	}

	return results, nil
}

func (client *Client) GetEmbeddingDimension() int {
	switch client.embeddingModel {
	case LargeEmbedding3:
		return DimensionModelLarge
	case SmallEmbedding3:
		return DimensionModelSmall
	default:
		return 0
	}
}

func (client *Client) GetModelId() string {
	return string(client.embeddingModel)
}
