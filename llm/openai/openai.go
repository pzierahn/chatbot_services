package openai

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
)

type Client struct {
	client         *openai.Client
	embeddingModel openai.EmbeddingModel
}

func New() (*Client, error) {
	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		return nil, fmt.Errorf("missing OPENAI_API_KEY")
	}

	return &Client{
		client:         openai.NewClient(token),
		embeddingModel: LargeEmbedding3,
	}, nil
}
