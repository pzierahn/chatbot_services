package voyageai

import (
	"errors"
	"os"
)

type Client struct {
	apiKey string
	model  string
}

func New(model string) (*Client, error) {
	apiKey := os.Getenv("VOYAGE_API_KEY")
	if apiKey == "" {
		return nil, errors.New("VOYAGEAI_API_KEY is not set")
	}

	return &Client{
		apiKey: apiKey,
		model:  model,
	}, nil
}
