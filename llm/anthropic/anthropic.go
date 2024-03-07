package anthropic

import (
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	"os"
)

type Client struct {
	usage  llm.Usage
	apiKey string
}

func New(usage llm.Usage) (*Client, error) {

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY env var is not set")
	}

	return &Client{
		usage:  usage,
		apiKey: apiKey,
	}, nil
}
