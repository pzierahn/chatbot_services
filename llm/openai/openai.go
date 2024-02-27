package openai

import (
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
	"os"
)

type Client struct {
	client *openai.Client
	usage  llm.Usage
}

func New(usage llm.Usage) (*Client, error) {
	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		return nil, fmt.Errorf("missing OPENAI_API_KEY")
	}

	return &Client{
		client: openai.NewClient(token),
		usage:  usage,
	}, nil
}
