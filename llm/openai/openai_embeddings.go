package openai

import (
	"context"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
)

const embeddingModel = openai.LargeEmbedding3

func (client *Client) CreateEmbeddings(ctx context.Context, req *llm.EmbeddingRequest) (*llm.EmbeddingResponse, error) {
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

	if !req.SkipTracking {
		client.usage.Track(ctx, llm.ModelUsage{
			UserId:       req.UserId,
			Model:        string(resp.Model),
			InputTokens:  resp.Usage.PromptTokens,
			OutputTokens: resp.Usage.CompletionTokens,
		})
	}

	return &llm.EmbeddingResponse{
		Data:   resp.Data[0].Embedding,
		Tokens: resp.Usage.PromptTokens,
	}, nil
}

func (client *Client) GetModelName() string {
	return string(embeddingModel)
}
