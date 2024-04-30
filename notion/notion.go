package notion

import (
	"fmt"
	"github.com/jomei/notionapi"
	"os"
)

type Client struct {
	api *notionapi.Client
}

func New() (*Client, error) {
	token := os.Getenv("NOTION_API_KEY")
	if token == "" {
		return nil, fmt.Errorf("missing NOTION_API_KEY")
	}

	return &Client{
		api: notionapi.NewClient(notionapi.Token(token)),
	}, nil
}
