package notion

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jomei/notionapi"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/chat"
	"github.com/pzierahn/chatbot_services/documents"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	pb.UnimplementedNotionServer
	chat      *chat.Service
	documents *documents.Service
	db        *pgxpool.Pool
	auth      auth.Service
}

func (client *Client) getAPIClient(ctx context.Context) (*notionapi.Client, error) {
	token, err := client.GetApiKey(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return notionapi.NewClient(notionapi.Token(token.Key)), nil
}

func New(chat *chat.Service, documents *documents.Service, db *pgxpool.Pool, auth auth.Service) (*Client, error) {
	return &Client{
		chat:      chat,
		documents: documents,
		db:        db,
		auth:      auth,
	}, nil
}
