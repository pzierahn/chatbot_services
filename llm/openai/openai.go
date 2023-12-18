package openai

import (
	"github.com/sashabaranov/go-openai"
	"os"
)

type Client struct {
	client *openai.Client
}

func New() *Client {
	token := os.Getenv("OPENAI_API_KEY")
	return &Client{
		client: openai.NewClient(token),
	}
}
