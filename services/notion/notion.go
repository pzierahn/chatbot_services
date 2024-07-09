package notion

import (
	"context"
	"github.com/jomei/notionapi"
	"github.com/pzierahn/chatbot_services/datastore"
	"github.com/pzierahn/chatbot_services/services/account"
	"github.com/pzierahn/chatbot_services/services/chat"
	"github.com/pzierahn/chatbot_services/services/documents"
	pb "github.com/pzierahn/chatbot_services/services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

type Client struct {
	pb.UnimplementedNotionServer
	Chat      *chat.Service
	Documents *documents.Service
	Database  *datastore.Service
	Auth      account.Verifier
	mux       sync.RWMutex
	Cache     map[string]string
}

func (client *Client) getCachedKey(userId string) (key string, ok bool) {
	client.mux.RLock()
	defer client.mux.RUnlock()
	key, ok = client.Cache[userId]
	return
}

func (client *Client) setCachedKey(userId, key string) {
	client.mux.Lock()
	defer client.mux.Unlock()
	client.Cache[userId] = key
}

func (client *Client) deleteCachedKey(userId string) {
	client.mux.Lock()
	defer client.mux.Unlock()
	delete(client.Cache, userId)
}

func (client *Client) getAPIClient(ctx context.Context) (*notionapi.Client, error) {
	token, err := client.GetAPIKey(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return notionapi.NewClient(notionapi.Token(token.Key)), nil
}
