package mistral

import (
	"fmt"
	"github.com/gage-technologies/mistral-go"
	"github.com/pzierahn/chatbot_services/llm"
	"os"
)

type Client struct {
	client *mistral.MistralClient
	usage  llm.Usage
}

func New(usage llm.Usage) (*Client, error) {

	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("MISTRAL_API_KEY env var is not set")
	}

	return &Client{
		client: mistral.NewMistralClientDefault(apiKey),
		usage:  usage,
	}, nil
}
