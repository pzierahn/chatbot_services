package notion

import (
	"context"
	"github.com/jomei/notionapi"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/chat"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/documents"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	pb.UnimplementedNotionServer
	Chat      *chat.Service
	Documents *documents.Service
	Database  *datastore.Service
	Auth      account.Verifier
}

func (client *Client) getAPIClient(ctx context.Context) (*notionapi.Client, error) {
	token, err := client.GetAPIKey(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return notionapi.NewClient(notionapi.Token(token.Key)), nil
}
