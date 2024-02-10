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
			Model: openai.LargeEmbedding3,
			Input: []string{req.Input},
			User:  req.UserId,
		},
	)
	if err != nil {
		return nil, err
	}

	client.trackUsage(ctx, llm.ModelUsage{
		UserId:           req.UserId,
		Model:            string(resp.Model),
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
	})

	return &llm.EmbeddingResponse{
		Data:   resp.Data[0].Embedding,
		Tokens: resp.Usage.PromptTokens,
	}, nil
}
