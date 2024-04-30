package notion

import (
	"fmt"
	"github.com/jomei/notionapi"
	pb "github.com/pzierahn/chatbot_services/proto"
	"os"
)

type Client struct {
	pb.UnimplementedNotionServer
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
