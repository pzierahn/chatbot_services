package openai

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
)

func (client *Client) CreateEmbeddings(ctx context.Context, req *llm.EmbeddingRequest) (*llm.EmbeddingResponse, error) {
	resp, err := client.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Model: "text-embedding-3-large",
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
