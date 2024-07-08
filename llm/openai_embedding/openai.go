package openai_embedding

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
)

type Client struct {
	client *openai.Client
	model  openai.EmbeddingModel
}

func New() (*Client, error) {
	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		return nil, fmt.Errorf("missing OPENAI_API_KEY")
	}

	return &Client{
		client: openai.NewClient(token),
	}, nil
}
