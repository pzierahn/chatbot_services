package openai_embedding

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
)

const embeddingModel = openai.LargeEmbedding3

func (client *Client) Create(ctx context.Context, req *llm.EmbeddingRequest) (*llm.EmbeddingResponse, error) {
	resp, err := client.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Model: embeddingModel,
			Input: []string{req.Input},
			User:  req.UserId,
		},
	)
	if err != nil {
		return nil, err
	}

	return &llm.EmbeddingResponse{
		Data:   resp.Data[0].Embedding,
		Tokens: resp.Usage.PromptTokens,
	}, nil
}

func (client *Client) GetEmbeddingModelName() string {
	return string(embeddingModel)
}
