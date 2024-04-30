package notion

import (
	"fmt"
	"github.com/jomei/notionapi"
	"github.com/pzierahn/chatbot_services/chat"
	"github.com/pzierahn/chatbot_services/documents"
	pb "github.com/pzierahn/chatbot_services/proto"
	"os"
)

type Client struct {
	pb.UnimplementedNotionServer
	api       *notionapi.Client
	chat      *chat.Service
	documents *documents.Service
}

func New(chat *chat.Service, documents *documents.Service) (*Client, error) {
	token := os.Getenv("NOTION_API_KEY")
	if token == "" {
		return nil, fmt.Errorf("missing NOTION_API_KEY")
	}

	return &Client{
		api:       notionapi.NewClient(notionapi.Token(token)),
		chat:      chat,
		documents: documents,
	}, nil
}
