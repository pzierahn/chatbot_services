package notion

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/chatbot_services/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

var mux sync.RWMutex
var apiKeyCache = make(map[string]string)

func (client *Client) SetApiKey(ctx context.Context, key *pb.NotionApiKey) (*emptypb.Empty, error) {

	userID, err := client.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	// Insert or update key in database
	_, err = client.db.Exec(ctx,
		`INSERT INTO notion_api_keys (user_id, api_key) VALUES ($1, $2) 
			   ON CONFLICT (user_id) DO UPDATE SET api_key = $2`, userID, key.Key)

	if err != nil {
		return nil, err
	}

	mux.Lock()
	apiKeyCache[userID] = key.Key
	mux.Unlock()

	return &emptypb.Empty{}, nil
}

func (client *Client) RemoveApiKey(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {

	userID, err := client.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	// Delete the key from the database
	_, err = client.db.Exec(ctx, `DELETE FROM notion_api_keys WHERE user_id = $1`, userID)

	if err != nil {
		return nil, err
	}

	mux.Lock()
	delete(apiKeyCache, userID)
	mux.Unlock()

	return &emptypb.Empty{}, nil
}

func (client *Client) GetApiKey(ctx context.Context, _ *emptypb.Empty) (*pb.NotionApiKey, error) {

	userID, err := client.auth.Verify(ctx)
	if err != nil {
		return nil, err
	}

	mux.RLock()
	key, ok := apiKeyCache[userID]
	mux.RUnlock()

	apiKey := &pb.NotionApiKey{}

	if ok {
		apiKey.Key = key
		return apiKey, nil
	}

	// Get the api key from the database
	err = client.db.QueryRow(ctx,
		`SELECT api_key FROM notion_api_keys WHERE user_id = $1`,
		userID).Scan(&apiKey.Key)
	if err != nil {
		return nil, fmt.Errorf("could not get api key: %v", err)
	}

	mux.Lock()
	apiKeyCache[userID] = apiKey.Key
	mux.Unlock()

	return apiKey, nil
}
